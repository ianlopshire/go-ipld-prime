package impls

import (
	ipld "github.com/ipld/go-ipld-prime/_rsrch/nodesolution"
)

var (
	//_ ipld.Node          = &anyNode{}
	_ ipld.NodeStyle   = Style__Any{}
	_ ipld.NodeBuilder = &anyBuilder{}
	//_ ipld.NodeAssembler = &anyAssembler{}
)

// anyNode is a union meant for alloc amortization; see anyAssembler.
// Note that anyBuilder doesn't use anyNode, because it's not aiming to amortize anything.
//
// REVIEW: if there's any point in keeping this around.  It's here for completeness,
// but not currently used anywhere in package, and also not currently exported.
// type anyNode struct {
// 	kind ipld.ReprKind
//
// 	plainMap
// 	plainList
// 	plainBool
// 	plainInt
// 	plainFloat
// 	plainString
// 	plainBytes
// 	plainLink
// }

// -- Node interface methods -->

// Unimplemented at present -- see "REVIEW" comment on anyNode.

// -- NodeStyle -->

type Style__Any struct{}

func (Style__Any) NewBuilder() ipld.NodeBuilder {
	return &anyBuilder{}
}

// -- NodeBuilder -->

// anyBuilder is a builder for any kind of node.
//
// anyBuilder is a little unusual in its internal workings:
// unlike most builders, it doesn't embed the corresponding assembler,
// nor will it end up using anyNode,
// but instead embeds a builder for each of the kinds it might contain.
// This is because we want a more granular return at the end:
// if we used anyNode, and returned a pointer to just the relevant part of it,
// we'd have all the extra bytes of anyNode still reachable in GC terms
// for as long as that handle to the interior of it remains live.
type anyBuilder struct {
	// kind is set on first interaction, and used to select which builder to delegate 'Build' to!
	// As soon as it's been set to a value other than zero (being "Invalid"), all other Assign/Begin calls will fail since something is already in progress.
	// May also be set to the magic value '99', which means "i dunno, I'm just carrying another node of unknown style".
	kind ipld.ReprKind

	// Only one of the following ends up being used...
	//  but we don't know in advance which one, so all are embeded here.
	//   This uses excessive space, but amortizes allocations, and all will be
	//    freed as soon as the builder is done.
	// Builders are only used for recursives;
	//  scalars are simple enough we just do them directly.
	// 'scalarNode' may also hold another Node of unknown style (possibly not even from this package),
	//  in which case this is indicated by 'kind==99'.

	mapBuilder  plainMap__Builder
	listBuilder plainList__Builder
	scalarNode  ipld.Node
}

func (nb *anyBuilder) Reset() {
	*nb = anyBuilder{}
}

func (nb *anyBuilder) BeginMap(sizeHint int) (ipld.MapNodeAssembler, error) {
	if nb.kind != ipld.ReprKind_Invalid {
		panic("misuse")
	}
	nb.kind = ipld.ReprKind_Map
	nb.mapBuilder.w = &plainMap{}
	return nb.mapBuilder.BeginMap(sizeHint)
}
func (nb *anyBuilder) BeginList(sizeHint int) (ipld.ListNodeAssembler, error) {
	if nb.kind != ipld.ReprKind_Invalid {
		panic("misuse")
	}
	nb.kind = ipld.ReprKind_List
	nb.listBuilder.w = &plainList{}
	return nb.listBuilder.BeginList(sizeHint)
}
func (nb *anyBuilder) AssignNull() error {
	if nb.kind != ipld.ReprKind_Invalid {
		panic("misuse")
	}
	nb.kind = ipld.ReprKind_Null
	return nil
}
func (nb *anyBuilder) AssignBool(v bool) error {
	if nb.kind != ipld.ReprKind_Invalid {
		panic("misuse")
	}
	nb.kind = ipld.ReprKind_Bool
	nb.scalarNode = Bool(v)
	return nil
}
func (nb *anyBuilder) AssignInt(v int) error {
	if nb.kind != ipld.ReprKind_Invalid {
		panic("misuse")
	}
	nb.kind = ipld.ReprKind_Int
	nb.scalarNode = Int(v)
	return nil
}
func (nb *anyBuilder) AssignFloat(v float64) error {
	if nb.kind != ipld.ReprKind_Invalid {
		panic("misuse")
	}
	nb.kind = ipld.ReprKind_Float
	nb.scalarNode = Float(v)
	return nil
}
func (nb *anyBuilder) AssignString(v string) error {
	if nb.kind != ipld.ReprKind_Invalid {
		panic("misuse")
	}
	nb.kind = ipld.ReprKind_String
	nb.scalarNode = String(v)
	return nil
}
func (nb *anyBuilder) AssignBytes(v []byte) error {
	if nb.kind != ipld.ReprKind_Invalid {
		panic("misuse")
	}
	nb.kind = ipld.ReprKind_Bytes
	nb.scalarNode = Bytes(v)
	return nil
}
func (nb *anyBuilder) AssignLink(v ipld.Link) error {
	if nb.kind != ipld.ReprKind_Invalid {
		panic("misuse")
	}
	nb.kind = ipld.ReprKind_Link
	nb.scalarNode = Link(v)
	return nil
}
func (nb *anyBuilder) AssignNode(v ipld.Node) error {
	if nb.kind != ipld.ReprKind_Invalid {
		panic("misuse")
	}
	nb.kind = 99
	nb.scalarNode = v
	return nil
}
func (anyBuilder) Style() ipld.NodeStyle {
	return Style__Any{}
}

func (nb *anyBuilder) Build() ipld.Node {
	switch nb.kind {
	case ipld.ReprKind_Invalid:
		panic("misuse")
	case ipld.ReprKind_Map:
		return nb.mapBuilder.Build()
	case ipld.ReprKind_List:
		return nb.listBuilder.Build()
	case ipld.ReprKind_Null:
		return ipld.Null
	case ipld.ReprKind_Bool:
		return nb.scalarNode
	case ipld.ReprKind_Int:
		return nb.scalarNode
	case ipld.ReprKind_Float:
		return nb.scalarNode
	case ipld.ReprKind_String:
		return nb.scalarNode
	case ipld.ReprKind_Bytes:
		return nb.scalarNode
	case ipld.ReprKind_Link:
		return nb.scalarNode
	case 99:
		return nb.scalarNode
	default:
		panic("unreachable")
	}
}

// -- NodeAssembler -->

// ... oddly enough, we seem to be able to put off implementing this
//  until we also implement something that goes full-hog on amortization
//   and actually has a slab of `anyNode`.  Which so far, nothing does.
//    See "REVIEW" comment on anyNode.
// type anyAssembler struct {
// 	w *anyNode
// }
