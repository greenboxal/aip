package forddbimpl

import (
	"context"
	"strconv"

	"github.com/ipld/go-ipld-prime"
)

type TreeDiff struct {
	Objects map[string]ObjectDiff `json:"objects"`
}

type ObjectDiff struct {
	Changes []ObjectChange `json:"changes"`
}

type ObjectChange struct {
	Path string    `json:"path"`
	Old  ipld.Node `json:"old"`
	New  ipld.Node `json:"new"`
}

func DiffTree(ctx context.Context, oldTree, newTree Tree, getter NodeGetter) (TreeDiff, error) {
	return TreeDiff{}, nil
}

// DiffObject compares two trees and returns a diff between them.
// The diff is a map of object IDs to a diff of the object.
// The diff of an object is a list of changes to the object.
// The changes are a path to the changed node and the old and new values.
// The path is a dot-separated list of keys to traverse the object.
// The old and new values are the values at the path.
func DiffObject(ctx context.Context, oldNode, newNode ipld.Node) (ObjectDiff, error) {
	diff := ObjectDiff{}

	err := diffNodes(oldNode, newNode, "", &diff.Changes)

	if err != nil {
		return diff, err
	}

	return diff, nil
}

func diffNodes(oldNode, newNode ipld.Node, path string, changes *[]ObjectChange) error {
	// Initialize oldMap and newMap for unordered comparison
	oldMap := make(map[string]ipld.Node)
	newMap := make(map[string]ipld.Node)

	if oldNode.Kind() == ipld.Kind_Map {
		oldIterator := oldNode.MapIterator()

		// Fill the oldMap
		for !oldIterator.Done() {
			oldKey, oldValue, err := oldIterator.Next()
			if err != nil {
				return err
			}
			str, err := oldKey.AsString()
			if err != nil {
				return err
			}
			oldMap[str] = oldValue
		}
	} else {
		oldIterator := oldNode.ListIterator()

		// Fill the oldMap
		for !oldIterator.Done() {
			oldIndex, oldValue, err := oldIterator.Next()
			if err != nil {
				return err
			}
			str := strconv.FormatInt(oldIndex, 10)
			oldMap[str] = oldValue
		}
	}

	if newNode.Kind() == ipld.Kind_Map {
		newIterator := newNode.MapIterator()

		// Fill the newMap
		for !newIterator.Done() {
			newKey, newValue, err := newIterator.Next()
			if err != nil {
				return err
			}
			str, err := newKey.AsString()
			if err != nil {
				return err
			}
			newMap[str] = newValue
		}
	} else {
		newIterator := newNode.ListIterator()

		// Fill the newMap
		for !newIterator.Done() {
			newIndex, newValue, err := newIterator.Next()
			if err != nil {
				return err
			}
			str := strconv.FormatInt(newIndex, 10)
			newMap[str] = newValue
		}
	}

	// Check keys in both oldMap and newMap
	for key, oldValue := range oldMap {
		newValue, ok := newMap[key]

		fullPath := path

		if path != "" {
			fullPath += "/"
		}

		fullPath += key

		if !ok {
			*changes = append(*changes, ObjectChange{
				Path: fullPath,
				Old:  oldValue,
				New:  nil,
			})
		} else {
			if oldValue.Kind() != newValue.Kind() {
				*changes = append(*changes, ObjectChange{
					Path: fullPath,
					Old:  oldValue,
					New:  newValue,
				})
			} else if newValue.Kind() == ipld.Kind_Map {
				err := diffNodes(oldValue, newValue, fullPath, changes)

				if err != nil {
					return err
				}
			} else if newValue.Kind() == ipld.Kind_List {
				err := diffNodes(oldValue, newValue, fullPath, changes)

				if err != nil {
					return err
				}
			} else if !mustCompareNodes(oldValue, newValue) {
				*changes = append(*changes, ObjectChange{
					Path: fullPath,
					Old:  oldValue,
					New:  newValue,
				})
			}
		}

		delete(newMap, key)
	}

	// Check the remaining keys in newMap (those only exist in newMap)
	for key, newValue := range newMap {
		fullPath := path

		if path != "" {
			fullPath += "/"
		}

		fullPath += key

		*changes = append(*changes, ObjectChange{
			Path: fullPath,
			Old:  nil,
			New:  newValue,
		})
	}

	return nil
}

func mustCompareNodes(value ipld.Node, value2 ipld.Node) bool {
	ok, err := compareNodes(value, value2)

	if err != nil {
		return false
	}

	return ok
}

func compareNodes(node1, node2 ipld.Node) (bool, error) {
	// If both nodes are nil, they are equal
	if node1 == nil && node2 == nil {
		return true, nil
	}

	// If one of the nodes is nil, they are not equal
	if node1 == nil || node2 == nil {
		return false, nil
	}

	// If node kinds are not equal, nodes are not equal
	if node1.Kind() != node2.Kind() {
		return false, nil
	}

	switch node1.Kind() {
	case ipld.Kind_Map:
		return compareMapNodes(node1, node2)
	case ipld.Kind_List:
		return compareListNodes(node1, node2)
	case ipld.Kind_String:
		string1, err1 := node1.AsString()
		string2, err2 := node2.AsString()
		if err1 != nil || err2 != nil {
			return false, nil
		}
		return string1 == string2, nil
	case ipld.Kind_Int:
		int1, err1 := node1.AsInt()
		int2, err2 := node2.AsInt()
		if err1 != nil || err2 != nil {
			return false, nil
		}
		return int1 == int2, nil
	case ipld.Kind_Float:
		float1, err1 := node1.AsFloat()
		float2, err2 := node2.AsFloat()
		if err1 != nil || err2 != nil {
			return false, nil
		}
		return float1 == float2, nil
	case ipld.Kind_Link:
		link1, err1 := node1.AsLink()
		link2, err2 := node2.AsLink()
		if err1 != nil || err2 != nil {
			return false, nil
		}
		return link1.String() == link2.String(), nil
	case ipld.Kind_Bool:
		bool1, err1 := node1.AsBool()
		bool2, err2 := node2.AsBool()
		if err1 != nil || err2 != nil {
			return false, nil
		}
		return bool1 == bool2, nil
	case ipld.Kind_Null:
		return true, nil // Both nodes are null, so they are equal
	case ipld.Kind_Bytes:
		bytes1, err1 := node1.AsBytes()
		bytes2, err2 := node2.AsBytes()
		if err1 != nil || err2 != nil {
			return false, nil
		}
		return string(bytes1) == string(bytes2), nil
	default:
		return false, nil
	}
}

func compareMapNodes(node1, node2 ipld.Node) (bool, error) {
	node1Iterator := node1.MapIterator()
	node2Iterator := node2.MapIterator()

	// Check the lengths of maps are the same
	if node1.Length() != node2.Length() {
		return false, nil
	}

	for !node1Iterator.Done() && !node2Iterator.Done() {
		key1, value1, err1 := node1Iterator.Next()
		key2, value2, err2 := node2Iterator.Next()

		if err1 != nil || err2 != nil || key1 != key2 {
			return false, nil
		}

		isEqual, err := compareNodes(value1, value2)
		if err != nil || !isEqual {
			return false, nil
		}
	}

	return true, nil
}

func compareListNodes(node1, node2 ipld.Node) (bool, error) {
	node1Iterator := node1.ListIterator()
	node2Iterator := node2.ListIterator()

	// Check the lengths of lists are the same
	if node1.Length() != node2.Length() {
		return false, nil
	}

	for !node1Iterator.Done() && !node2Iterator.Done() {
		_, value1, err1 := node1Iterator.Next()
		_, value2, err2 := node2Iterator.Next()

		if err1 != nil || err2 != nil {
			return false, nil
		}

		isEqual, err := compareNodes(value1, value2)
		if err != nil || !isEqual {
			return false, nil
		}
	}

	return true, nil
}
