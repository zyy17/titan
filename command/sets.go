package command

import (
	"bytes"
	"errors"
	"strconv"

	"github.com/meitu/titan/db"
)

// SAdd adds the specified members to the set stored at key
func SAdd(ctx *Context, txn *db.Transaction) (OnCommit, error) {
	key := []byte(ctx.Args[0])

	members := make([][]byte, len(ctx.Args[1:]))
	for i, member := range ctx.Args[1:] {
		members[i] = []byte(member)
	}
	set, err := txn.Set(key)
	if err != nil {
		return nil, errors.New("ERR " + err.Error())
	}
	added, err := set.SAdd(members...)
	if err != nil {
		return nil, errors.New("ERR " + err.Error())
	}
	return Integer(ctx.Out, added), nil
}

// SMembers returns all the members of the set value stored at key
func SMembers(ctx *Context, txn *db.Transaction) (OnCommit, error) {
	key := []byte(ctx.Args[0])

	set, err := txn.Set(key)
	if err != nil {
		return nil, errors.New("ERR " + err.Error())
	}

	members, err := set.SMembers()
	if err != nil {
		return nil, errors.New("ERR " + err.Error())
	}
	return BytesArray(ctx.Out, members), nil
}

// SCard returns the set cardinality (number of elements) of the set stored at key
func SCard(ctx *Context, txn *db.Transaction) (OnCommit, error) {
	key := []byte(ctx.Args[0])

	set, err := txn.Set(key)
	if err != nil {
		return nil, errors.New("ERR " + err.Error())
	}
	count, err := set.SCard()
	if err != nil {
		return nil, errors.New("ERR " + err.Error())
	}
	return Integer(ctx.Out, int64(count)), nil
}

// SIsmember returns if member is a member of the set stored at key
func SIsmember(ctx *Context, txn *db.Transaction) (OnCommit, error) {
	key := []byte(ctx.Args[0])
	member := []byte(ctx.Args[1])
	set, err := txn.Set(key)
	if err != nil {
		return nil, errors.New("ERR " + err.Error())
	}
	count, err := set.SIsmember(member)
	if err != nil {
		return nil, errors.New("ERR " + err.Error())
	}
	return Integer(ctx.Out, int64(count)), nil

}

// SPop removes and returns one or more random elements from the set value store at key
func SPop(ctx *Context, txn *db.Transaction) (OnCommit, error) {
	var count int
	var err error
	var members [][]byte
	var set *db.Set
	key := []byte(ctx.Args[0])

	if len(ctx.Args) == 2 {
		count, err = strconv.Atoi(ctx.Args[1])
		if err != nil {
			return nil, errors.New("ERR " + err.Error())
		}
	}
	set, err = txn.Set(key)
	if err != nil {
		return nil, errors.New("ERR " + err.Error())
	}
	members, err = set.SPop(int64(count))
	if err != nil {
		return nil, errors.New("ERR " + err.Error())
	}
	return BytesArray(ctx.Out, members), nil
}

// SRem removes the specified members from the set stored at key
func SRem(ctx *Context, txn *db.Transaction) (OnCommit, error) {
	var members [][]byte
	key := []byte(ctx.Args[0])
	for _, member := range ctx.Args[1:] {
		members = append(members, []byte(member))
	}
	set, err := txn.Set(key)
	if err != nil {
		return nil, errors.New("ERR " + err.Error())
	}
	count, err := set.SRem(members)
	if err != nil {
		return nil, errors.New("ERR " + err.Error())
	}
	return Integer(ctx.Out, int64(count)), nil
}

// SMove movies member from the set at source to the set at destination
func SMove(ctx *Context, txn *db.Transaction) (OnCommit, error) {
	member := make([]byte, 0, len(ctx.Args[2]))
	key := []byte(ctx.Args[0])
	destkey := []byte(ctx.Args[1])
	member = []byte(ctx.Args[2])

	set, err := txn.Set(key)
	if err != nil {
		return nil, errors.New("ERR " + err.Error())
	}
	count, err := set.SMove(destkey, member)
	if err != nil {
		return nil, errors.New("ERR " + err.Error())
	}
	return Integer(ctx.Out, int64(count)), nil
}

// SUnion returns the members of the set resulting from the union of all the given sets.
func SUnion(ctx *Context, txn *db.Transaction) (OnCommit, error) {
	var members [][]byte
	mambermap := make(map[string]int)
	keys := make([][]byte, len(ctx.Args))
	for i, key := range ctx.Args {
		keys[i] = []byte(key)

		set, err := txn.Set(keys[i])
		if err != nil {
			return nil, errors.New("ERR " + err.Error())
		}
		if !set.Exists() {
			continue
		}
		if n, _ := set.SCard(); n == 0 {
			continue
		}
		ms, err := set.SMembers()
		if err != nil {
			return nil, errors.New("ERR " + err.Error())
		}
		for n := range ms {
			mambermap[string(ms[n])] = 1
		}

	}
	for k, _ := range mambermap {
		members = append(members, []byte(k))
	}

	return BytesArray(ctx.Out, members), nil
}

// SInter returns the members of the set resulting from the intersection of all the given sets.
func SInter(ctx *Context, txn *db.Transaction) (OnCommit, error) {
	var members [][]byte
	keys := make([][]byte, len(ctx.Args[1:]))
	key := []byte(ctx.Args[0])

	for i, key := range ctx.Args[1:] {
		keys[i] = []byte(key)
	}
	set, err := txn.Set(key)
	if err != nil {
		return nil, errors.New("ERR " + err.Error())
	}

	if !set.Exists() {
		return nil, nil
	}
	if n, _ := set.SCard(); n == 0 {
		return nil, nil
	}
	members, err = set.SMembers()
	if err != nil {
		return nil, errors.New("ERR " + err.Error())
	}

	for i := range keys {
		set, err := txn.Set(keys[i])
		if err != nil {
			return nil, errors.New("ERR " + err.Error())
		}
		if !set.Exists() {
			return nil, nil
		}
		if n, _ := set.SCard(); n == 0 {
			return nil, nil
		}
		ms, err := set.SMembers()
		if err != nil {
			return nil, errors.New("ERR " + err.Error())
		}
		members = sliceInter(members, ms)

	}
	return BytesArray(ctx.Out, members), nil
}

// SDiff returns the members of the set resulting from the difference between the first set and all the successive sets.
func SDiff(ctx *Context, txn *db.Transaction) (OnCommit, error) {
	var keys [][]byte
	var members [][]byte

	key := []byte(ctx.Args[0])
	set, err := txn.Set(key)
	if err != nil {
		return nil, errors.New("ERR " + err.Error())
	}
	if !set.Exists() {
		return nil, nil
	}
	if n, _ := set.SCard(); n == 0 {
		return nil, nil
	}

	members, err = set.SMembers()
	if err != nil {
		return nil, errors.New("ERR " + err.Error())
	}

	for _, key := range ctx.Args[1:] {
		keys = append(keys, []byte(key))
	}

	for i := range keys {
		set, err := txn.Set(keys[i])
		if err != nil {
			return nil, errors.New("ERR " + err.Error())
		}
		if !set.Exists() {
			continue
		}
		if n, _ := set.SCard(); n == 0 {
			continue
		}
		ms, err := set.SMembers()
		if err != nil {
			return nil, errors.New("ERR " + err.Error())
		}
		members = sliceDiff(members, ms)
	}
	return BytesArray(ctx.Out, members), nil
}

// InSliceInter checks given interface in interface slice.
func inSliceInter(v []byte, sl [][]byte) bool {
	for _, vv := range sl {
		if bytes.Equal(vv, v) {
			return true
		}
	}
	return false
}

// SliceIntersect returns slice that are present in all the slice1 and slice2.
func sliceInter(slice1, slice2 [][]byte) (interslice [][]byte) {
	for _, v := range slice1 {
		if inSliceInter(v, slice2) {
			interslice = append(interslice, v)
		}
	}
	return
}

// InSliceDiff checks given interface in interface slice.
func inSliceDiff(v []byte, sl [][]byte) bool {
	for _, vv := range sl {
		if bytes.Equal(vv, v) {
			return false
		}
	}
	return true
}

// SliceIntersect returns all slices in slice1 that are not present in slice2.
func sliceDiff(slice1, slice2 [][]byte) [][]byte {
	var diffslice [][]byte
	for _, v := range slice1 {
		if inSliceDiff(v, slice2) {
			diffslice = append(diffslice, v)
		}
	}
	return diffslice
}
