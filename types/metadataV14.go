package types

import (
	"fmt"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v3/scale"
)

// Based on:
// https://github.com/polkadot-js/api/blob/48ef04b8ca21dc4bd06442775d9b7585c75d1253/packages/types/src/interfaces/metadata/v14.ts

type MetadataV14 struct {
	Lookup    PortableTypeV14
	Pallets   []PalletMetadataV14
	Extrinsic ExtrinsicMetadataV14
	Type      Si1LookupTypeId
}

func (m *MetadataV14) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Lookup)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Pallets)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Extrinsic)
	if err != nil {
		return err
	}

	return decoder.Decode(&m.Type)
}

func (m MetadataV14) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.Pallets)
	if err != nil {
		return err
	}
	return encoder.Encode(m.Extrinsic)
}

func (m *MetadataV14) FindCallIndex(call string) (CallIndex, error) {
	s := strings.Split(call, ".")
	for _, mod := range m.Pallets {
		if !mod.HasCalls {
			continue
		}
		if string(mod.Name) != s[0] {
			continue
		}
		for ci, f := range mod.Constants {
			if string(f.Name) == s[1] {
				return CallIndex{mod.Index, uint8(ci)}, nil
			}
		}
		return CallIndex{}, fmt.Errorf("method %v not found within module %v for call %v", s[1], mod.Name, call)
	}
	return CallIndex{}, fmt.Errorf("module %v not found in metadata for call %v", s[0], call)
}

func (m *MetadataV14) FindEventNamesForEventID(eventID EventID) (Text, Text, error) {
	for _, mod := range m.Pallets {
		if !mod.HasEvents {
			continue
		}
		if mod.Index != eventID[0] {
			continue
		}
		//if int(eventID[1]) >= len(mod.Events.Type) {
		//	return "", "", fmt.Errorf("event index %v for module %v out of range", eventID[1], mod.Name)
		//}
		return mod.Name, mod.Constants[eventID[1]].Name, nil
	}
	return "", "", fmt.Errorf("module index %v out of range", eventID[0])
}

func (m *MetadataV14) FindStorageEntryMetadata(module string, fn string) (StorageEntryMetadata, error) {
	for _, mod := range m.Pallets {
		if !mod.HasStorage {
			continue
		}
		if string(mod.Storage.Prefix) != module {
			continue
		}
		for _, s := range mod.Storage.Items {
			if string(s.Name) != fn {
				continue
			}
			return nil, nil
		}
		return nil, fmt.Errorf("storage %v not found within module %v", fn, module)
	}
	return nil, fmt.Errorf("module %v not found in metadata", module)
}

func (m *MetadataV14) FindConstantValue(module Text, constant Text) ([]byte, error) {
	for _, mod := range m.Pallets {
		if mod.Name == module {
			value, err := mod.FindConstantValue(constant)
			if err == nil {
				return value, nil
			}
		}
	}
	return nil, fmt.Errorf("could not find constant %s.%s", module, constant)
}

func (m *MetadataV14) ExistsModuleMetadata(module string) bool {
	for _, mod := range m.Pallets {
		if string(mod.Name) == module {
			return true
		}
	}
	return false
}

type SignedExtensionMetadataV14 struct {
	Identifier       Text
	Type             Si1LookupTypeId
	AdditionalSigned Si1LookupTypeId
}

type ExtrinsicMetadataV14 struct {
	Type             Si1LookupTypeId
	Version          uint8
	SignedExtensions []SignedExtensionMetadataV14
}

func (e *ExtrinsicMetadataV14) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&e.Type)
	if err != nil {
		return err
	}

	err = decoder.Decode(&e.Version)
	if err != nil {
		return err
	}

	return decoder.Decode(&e.SignedExtensions)
}

func (e ExtrinsicMetadataV14) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(&e.Type)
	if err != nil {
		return err
	}

	err = encoder.Encode(e.Version)
	if err != nil {
		return err
	}

	return encoder.Encode(e.SignedExtensions)
}

type PortableTypeV14 struct {
	//todo implement SiType/Si1Type https://github.com/polkadot-js/api/blob/2801088b0a05e6bc505c2c449f1eddb31b15587d/packages/types/src/interfaces/scaleInfo/v1.ts#L25
	Types []SiType
}

type SiType struct {
	Id Si1LookupTypeId
	Type Si1Type
}

type Si1Type struct {
	Path   Si1Path
	Params []Si1TypeParameter
	Def    Si1TypeDef
	Docs   []Text
}

type Si1Field struct {
	HasName     bool
	Name        Text
	Type        Si1LookupTypeId
	HasTypeName bool
	TypeName    Text
	Docs        []Text
}

func (s *Si1Field) Decode(decoder scale.Decoder) error {
	err := decoder.DecodeOption(&s.HasName, &s.Name)
	if err != nil {
		return err
	}

	err = decoder.Decode(&s.Type)
	if err != nil {
		return err
	}

	err = decoder.DecodeOption(&s.HasTypeName, &s.TypeName)
	if err != nil {
		return err
	}

	return decoder.Decode(&s.Docs)
}

type Si1TypeDefComposite struct {
	Fields []Si1Field
}

type Si1Variant struct {
	Name   Text
	Fields []Si1Field
	Index  U8
	Docs   []Text
}

func (s *Si1Variant) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&s.Name)
	if err != nil {
		return err
	}

	err = decoder.Decode(&s.Fields)
	if err != nil {
		return err
	}

	err = decoder.Decode(&s.Index)
	if err != nil {
		return err
	}

	return decoder.Decode(&s.Docs)
}

type Si1TypeDefVariant struct {
	Variants []Si1Variant
}

type Si1TypeDefSequence struct {
	Type Si1LookupTypeId
}

type Si1TypeDefArray struct {
	Len  U32
	Type Si1LookupTypeId
}

type Si1TypeDefTuple []Si1LookupTypeId

type Si1TypeDefPrimitive Si0TypeDefPrimitive

type Si0TypeDefPrimitive struct {
	IsBool bool
	IsChar bool
	IsStr  bool
	IsU8   bool
	IsU16  bool
	IsU32  bool
	IsU64  bool
	IsU128 bool
	IsU256 bool
	IsI8   bool
	IsI16  bool
	IsI32  bool
	IsI64  bool
	IsI128 bool
	IsI256 bool
}

func (s *Si0TypeDefPrimitive) Decode(decoder scale.Decoder) error {
	var t uint8
	err := decoder.Decode(&t)
	if err != nil {
		return err
	}

	switch t {
	case 0:
		s.IsBool = true
	case 1:
		s.IsChar = true
	case 2:
		s.IsStr = true
	case 3:
		s.IsU8 = true
	case 4:
		s.IsU16 = true
	case 5:
		s.IsU32 = true
	case 6:
		s.IsU64 = true
	case 7:
		s.IsU128 = true
	case 8:
		s.IsU256 = true
	case 9:
		s.IsI8 = true
	case 10:
		s.IsI16 = true
	case 11:
		s.IsI32 = true
	case 12:
		s.IsI64 = true
	case 13:
		s.IsI128 = true
	case 14:
		s.IsI256 = true

	default:
		return fmt.Errorf("received unexpected type %v", t)
	}
	return nil
}

type Si1TypeDefCompact struct {
	Type Si1LookupTypeId
}

type Si1TypeDefBitSequence struct {
	BitStoreType Si1LookupTypeId
	BitOrderType Si1LookupTypeId
}

type Si1TypeDef struct {
	IsComposite          bool
	AsComposite          Si1TypeDefComposite
	IsVariant            bool
	AsVariant            Si1TypeDefVariant
	IsSequence           bool
	AsSequence           Si1TypeDefSequence
	IsArray              bool
	AsArray              Si1TypeDefArray
	IsTuple              bool
	AsTuple              Si1TypeDefTuple
	IsPrimitive          bool
	AsPrimitive          Si0TypeDefPrimitive
	IsCompact            bool
	AsCompact            Si1TypeDefCompact
	IsBitSequence        bool
	AsBitSequence        Si1TypeDefBitSequence
	IsHistoricMetaCompat bool
	AsHistoricMetaCompat Type
}

func (s *Si1TypeDef) Decode(decoder scale.Decoder) error {
	var t uint8
	err := decoder.Decode(&t)
	if err != nil {
		return err
	}

	switch t {
	case 0:
		s.IsComposite = true
		err = decoder.Decode(&s.AsComposite)
		if err != nil {
			return err
		}
	case 1:
		s.IsVariant = true
		err = decoder.Decode(&s.AsVariant)
		if err != nil {
			return err
		}
	case 2:
		s.IsSequence = true
		err = decoder.Decode(&s.AsSequence)
		if err != nil {
			return err
		}
	case 3:
		s.IsArray = true
		err = decoder.Decode(&s.AsArray)
		if err != nil {
			return err
		}
	case 4:
		s.IsTuple = true
		err = decoder.Decode(&s.AsTuple)
		if err != nil {
			return err
		}
	case 5:
		s.IsPrimitive = true
		err = decoder.Decode(&s.AsPrimitive)
		if err != nil {
			return err
		}
	case 6:
		s.IsCompact = true
		err = decoder.Decode(&s.AsCompact)
		if err != nil {
			return err
		}
	case 7:
		s.IsBitSequence = true
		err = decoder.Decode(&s.AsBitSequence)
		if err != nil {
			return err
		}
	case 8:
		s.IsHistoricMetaCompat = true
		err = decoder.Decode(&s.AsHistoricMetaCompat)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("received unexpected type %v", t)
	}
	return nil
}

type Si1TypeParameter struct {
	Name    Text
	HasType bool
	Type    Si1LookupTypeId
}

func (s *Si1TypeParameter) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&s.Name)
	if err != nil {
		return err
	}

	err = decoder.DecodeOption(&s.HasType, &s.Type)
	if err != nil {
		return err
	}


	return err
}

//type Si1LookupTypeId Si1LookupTypeId

type Si1LookupTypeId struct {
	U32
}

//func NewSi1LookupTypeId(i big.Int) Si1LookupTypeId {
//	return Si1LookupTypeId{&i}
//}

func (s *Si1LookupTypeId) Decode(decoder scale.Decoder) error {
	b, err := decoder.DecodeUintCompact()
	if err != nil {
		return err
	}

	*s = Si1LookupTypeId{U32(b.Uint64())}
	return nil
}

type Si1Path Si0Path

type Si0Path []Text

type PalletCallMetadataV14 struct {
	Type Si1LookupTypeId
}

type PalletEventMetadataV14 struct {
	Type Si1LookupTypeId
}

type PalletConstantMetadataV14 struct {
	Name  Text
	Type  Si1LookupTypeId
	Value Bytes
	Docs  []Text
}

type PalletErrorMetadataV14 struct {
	Type Si1LookupTypeId
}

type PalletMetadataV14 struct {
	Name       Text
	HasStorage bool
	Storage    PalletStorageMetadataV14
	HasCalls   bool
	Calls      PalletCallMetadataV14
	HasEvents  bool
	Events     PalletEventMetadataV14
	Constants  []PalletConstantMetadataV14
	HasErrors  bool
	Errors     PalletErrorMetadataV14
	Index      uint8
}

func (m *PalletMetadataV14) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Name)
	if err != nil {
		return err
	}

	err = decoder.DecodeOption(&m.HasStorage, &m.Storage)
	if err != nil {
		return err
	}

	err = decoder.DecodeOption(&m.HasCalls, &m.Calls)
	if err != nil {
		return err
	}

	err = decoder.DecodeOption(&m.HasEvents, &m.Events)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Constants)
	if err != nil {
		return err
	}

	err = decoder.DecodeOption(&m.HasErrors, &m.Errors)
	if err != nil {
		return err
	}

	return decoder.Decode(&m.Index)
}

func (m PalletMetadataV14) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.Name)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.HasStorage)
	if err != nil {
		return err
	}

	if m.HasStorage {
		err = encoder.Encode(m.Storage)
		if err != nil {
			return err
		}
	}

	err = encoder.Encode(m.HasCalls)
	if err != nil {
		return err
	}

	if m.HasCalls {
		err = encoder.Encode(m.Calls)
		if err != nil {
			return err
		}
	}

	err = encoder.Encode(m.HasEvents)
	if err != nil {
		return err
	}

	if m.HasEvents {
		err = encoder.Encode(m.Events)
		if err != nil {
			return err
		}
	}

	err = encoder.Encode(m.Constants)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Errors)
	if err != nil {
		return err
	}

	return encoder.Encode(m.Index)
}

func (m *PalletMetadataV14) FindConstantValue(constant Text) ([]byte, error) {
	for _, cons := range m.Constants {
		if cons.Name == constant {
			return cons.Value, nil
		}
	}
	return nil, fmt.Errorf("could not find constant %s", constant)
}

type PalletStorageMetadataV14 struct {
	Prefix Text
	Items  []StorageEntryMetadataV14
}

type StorageEntryModifierV14 StorageFunctionModifierV0

type StorageEntryMetadataV14 struct {
	Name     Text
	Modifier StorageFunctionModifierV0
	Type     StorageEntryTypeV14
	Fallback Bytes
	Docs     []Text
}

func (s StorageEntryMetadataV14) IsPlain() bool {
	return s.Type.IsType
}

func (s StorageEntryMetadataV14) IsMap() bool {
	return s.Type.IsMap
}

//func (s StorageEntryMetadataV14) Hasher() (hash.Hash, error) {
//	if s.Type.IsMap {
//		return s.Type.AsMap.Hasher.HashFunc()
//	}
//	return xxhash.New128(nil), nil
//}
//
//func (s StorageEntryMetadataV14) Hasher2() (hash.Hash, error) {
//	if !s.Type.IsDoubleMap {
//		return nil, fmt.Errorf("only DoubleMaps have a Hasher2")
//	}
//	return s.Type.AsDoubleMap.Key2Hasher.HashFunc()
//}
//
//func (s StorageEntryMetadataV14) Hashers() ([]hash.Hash, error) {
//	if !s.Type.IsNMap {
//		return nil, fmt.Errorf("only NMaps have Hashers")
//	}
//
//	hashers := make([]hash.Hash, len(s.Type.AsNMap.Hashers))
//	for i, hasher := range s.Type.AsNMap.Hashers {
//		hasherFn, err := hasher.HashFunc()
//		if err != nil {
//			return nil, err
//		}
//		hashers[i] = hasherFn
//	}
//	return hashers, nil
//}

type StorageHasherV14 StorageHasherV10

type MapTypeV14 struct {
	Hasher []StorageHasherV10
	Key    Si1LookupTypeId
	Value  Si1LookupTypeId
}

type StorageEntryTypeV14 struct {
	IsType      bool
	AsType      Si1LookupTypeId // 0
	IsMap       bool
	AsMap       MapTypeV14 // 2
}

func (s *StorageEntryTypeV14) Decode(decoder scale.Decoder) error {
	var t uint8
	err := decoder.Decode(&t)
	if err != nil {
		return err
	}

	switch t {
	case 0:
		s.IsType = true
		err = decoder.Decode(&s.AsType)
		if err != nil {
			return err
		}
	case 1:
		s.IsMap = true
		err = decoder.Decode(&s.AsMap)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("received unexpected type %v", t)
	}
	return nil
}

func (s StorageEntryTypeV14) Encode(encoder scale.Encoder) error {
	switch {
	case s.IsType:
		err := encoder.PushByte(0)
		if err != nil {
			return err
		}
		err = encoder.Encode(s.AsType)
		if err != nil {
			return err
		}
	case s.IsMap:
		err := encoder.PushByte(1)
		if err != nil {
			return err
		}
		err = encoder.Encode(s.AsMap)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("expected to be either type, map, double map or nmap but none was set: %v", s)
	}
	return nil
}

type NMapTypeV14 struct {
	Keys    []Type
	Hashers []StorageHasherV10
	Value   Type
}
