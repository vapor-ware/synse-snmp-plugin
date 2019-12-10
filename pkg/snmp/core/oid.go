package core

import (
	"fmt"
	"strconv"
	"strings"
)

// Oid is an SNMP OID.
type Oid struct {
	// The string representation of the Oid.
	ToString string
	// The slice representation of the Oid. Example [1 3 6 1 2 1 33 1 2 1 0]
	// There is technically no limit on the size of the integer in an ASN1 OBJECT
	// IDENTIFIER node. Here we handle max uint64.
	ToSlice []uint64
}

// NewOid creates an SNMP OID from a string.
func NewOid(oid string) (result *Oid, err error) {

	// Remove leading .
	oid = strings.TrimPrefix(oid, ".")

	if len(oid) == 0 {
		return nil, fmt.Errorf("empty oid")
	}

	split := strings.Split(oid, ".")
	var slice []uint64

	var segment uint64
	for i := 0; i < len(split); i++ {

		// Convert the segment to uint64.
		segment, err = strconv.ParseUint(split[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"split[%d] %v is not a uint64, %T", i, split[i], split[i])
		}

		// Add to the oid.
		slice = append(slice, segment)
	}

	result = &Oid{
		ToString: oid,
		ToSlice:  slice,
	}
	return result, nil
}

// NewOidFromSlice creates an SNMP OID from a slice.
// Example slice: [1 3 6 1 2 1 33 1 2 1 0] (UPS_MIB upsBatteryStatus)
func NewOidFromSlice(oid []uint64) (result *Oid, err error) {

	s := "." + strings.Trim(strings.Replace(fmt.Sprint(oid), " ", ".", -1), "[]")

	result = &Oid{
		ToString: s,
		ToSlice:  oid,
	}
	// Not currently an error here, but here in case we need one on the interface.
	return result, nil
}

// OidSlice is a slice of pointers to Oid.
type OidSlice []*Oid

// NewOids creates a slice of *Oid from a string slice.
func NewOids(oids []string) (result OidSlice, err error) {
	var oid *Oid
	for i := 0; i < len(oids); i++ {
		oid, err = NewOid(oids[i])
		if err != nil {
			return nil, err
		}
		result = append(result, oid)
	}
	return result, nil
}

// oidTrieNode contains a segment of an OID (The part between the periods) and
// pointers to children.
type oidTrieNode struct {
	OidSegment uint64         // The portion of the OID between periods.
	Children   []*oidTrieNode // N way trie. Pointers to child nodes.
}

// newOidTrieNode creates an oidTrieNode with no children.
func newOidTrieNode(oidSegment uint64) (node *oidTrieNode) {
	return &oidTrieNode{
		OidSegment: oidSegment,
		Children:   nil,
	}
}

// Dump dumps an oidTrieNode to the console.
func (oidTrieNode *oidTrieNode) Dump() (err error) {
	if oidTrieNode == nil {
		return fmt.Errorf("oidTrieNode is nil")
	}
	childCount := 0
	if oidTrieNode.Children != nil {
		childCount = len(oidTrieNode.Children)
	}
	fmt.Printf("oidTrieNode: segment %v, child count %v\n",
		oidTrieNode.OidSegment, childCount)
	return nil
}

// OidTrie is an n-way trie of Oids. Each SNMP OID is stored similarly to words
// in a trie.
type OidTrie struct {
	Head *oidTrieNode
}

// NewOidTrie creates an OidTrie with a head with segment zero.
// If oids are given as strings, they are inserted.
func NewOidTrie(oids *[]string) (trie *OidTrie, err error) {
	oidTrie := &OidTrie{Head: newOidTrieNode(0)}

	if oids == nil {
		return oidTrie, nil
	}

	var newOids OidSlice
	newOids, err = NewOids(*oids)
	if err != nil {
		return nil, err
	}

	// Insert each oid to the trie.
	for i := 0; i < len(newOids); i++ {
		err = oidTrie.Insert(newOids[i])
		if err != nil {
			return nil, err
		}
	}
	return oidTrie, nil
}

// Insert an Oid into the trie.
func (oidTrie *OidTrie) Insert(oid *Oid) (err error) {
	if oid == nil {
		return fmt.Errorf("oid is nil")
	}

	current := oidTrie.Head // Start at the head.
	for i := 0; i < len(oid.ToSlice); i++ {
		// Insert each segment and traverse down the trie.
		current, err = insertSegment(current, oid.ToSlice[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// insertSegment inserts one OID segment (the part between the periods) into
// the OidTrie. Returns the node at the position of insertion. If the node
// exists no duplicate insertion is made, but the new position is returned.
func insertSegment(
	current *oidTrieNode, segment uint64) (inserted *oidTrieNode, err error) {
	if current == nil {
		return nil, fmt.Errorf("current is nil")
	}
	// No children case.
	if current.Children == nil {
		current.Children = []*oidTrieNode{newOidTrieNode(segment)}
		return current.Children[0], nil
	}

	// Search for segment. This is linear search which is improvable later.
	// The trie will generally be sparse.
	for i := 0; i < len(current.Children); i++ {
		if current.Children[i].OidSegment == segment {
			// Node already inserted. Return new current position.
			return current.Children[i], nil
		}
		if current.Children[i].OidSegment > segment {
			// Node is not in there. Shift the slice and insert at position i.
			newNode := newOidTrieNode(segment)
			current.Children = append(current.Children, newOidTrieNode(0))
			copy(current.Children[i+1:], current.Children[i:])
			current.Children[i] = newNode
			return current.Children[i], nil
		}
	}
	// Not found. We have the largest segment at this level. Insert at end.
	newNode := newOidTrieNode(segment)
	current.Children = append(current.Children, newNode)
	return newNode, nil
}

// TrieIterator creates an iterator that will traverse the trie in order via
// Next().
type TrieIterator struct {
	Head *oidTrieNode // Pointer to the head.
	// Below are stacks.
	Accumulator    []uint64       // Accumulation of the OID.
	NodeFrom       []*oidTrieNode // Which node we came from.
	NextChildIndex []int          // Which child pointer we traverse to next.
}

// NewTrieIterator initializes the OidTrie iterator.
func NewTrieIterator(oidTrie *OidTrie) (trieIterator *TrieIterator) {
	return &TrieIterator{
		Head:           oidTrie.Head,                 // Start at head.
		NodeFrom:       []*oidTrieNode{oidTrie.Head}, // Start at head.
		NextChildIndex: []int{-1},                    // Have not yet started.
	}
}

// Next traverses the Oid trie with a TrieIterator in order. Returns a pointer
// to one Oid at a time.
func (trieIterator *TrieIterator) Next() (oid *Oid, done bool, err error) {
	if trieIterator == nil {
		return nil, true, fmt.Errorf("trieIterator is nil")
	}

	var index int
	for {
		// currentNode is the last item in the NodeFrom stack.
		currentNode := trieIterator.NodeFrom[len(trieIterator.NodeFrom)-1]
		tailIndex := len(trieIterator.NextChildIndex) - 1
		index = trieIterator.NextChildIndex[tailIndex] // index we are processing this iteration.
		trieIterator.NextChildIndex[tailIndex]++       // index to process on the next iteration.

		if index == -1 {
			if len(currentNode.Children) == 0 {
				// Return an Oid from the accumulator since we are on the leaf.
				oid, err = NewOidFromSlice(trieIterator.Accumulator)
				if err != nil {
					return nil, true, nil
				}
				return oid, false, nil
			}
			continue
		}

		if index == len(currentNode.Children) {
			// End of current node.
			if currentNode == trieIterator.Head {
				return nil, true, nil // Done.
			}
			// Backtrack (Pop) currentNode, next child index, accumulator.
			trieIterator.NodeFrom = trieIterator.NodeFrom[:len(trieIterator.NodeFrom)-1]
			trieIterator.NextChildIndex = trieIterator.NextChildIndex[:len(trieIterator.NextChildIndex)-1]
			trieIterator.Accumulator = trieIterator.Accumulator[:len(trieIterator.Accumulator)-1]
			continue
		}

		if currentNode != nil && index < len(currentNode.Children) {
			// Push.
			trieIterator.Accumulator = append(trieIterator.Accumulator, currentNode.Children[index].OidSegment)
			trieIterator.NodeFrom = append(trieIterator.NodeFrom, currentNode.Children[index])
			trieIterator.NextChildIndex = append(trieIterator.NextChildIndex, -1)
		}
	}
}

// Sort returns a slice of SNMP OIDs from the trie in SNMP sort order.
func (oidTrie *OidTrie) Sort() (result []*Oid, err error) {

	var oid *Oid
	done := false
	trieIterator := NewTrieIterator(oidTrie)

	for !done {
		oid, done, err = trieIterator.Next()
		if err != nil {
			return nil, err
		}
		if done {
			break
		}

		// Add to result.
		result = append(result, oid)
	}

	return result, nil
}
