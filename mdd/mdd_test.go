package mdd

import (
	"encoding/xml"
	"strconv"
	"testing"

	"github.com/matrixxsoftware/go-mdd/dictionary"
	"github.com/matrixxsoftware/go-mdd/mdd/field"
	"github.com/stretchr/testify/assert"
)

func TestGetContainer(t *testing.T) {

	mdc := Containers{
		Containers: []Container{
			{
				Header: Header{
					Version:       1,
					TotalField:    5,
					Depth:         0,
					Key:           101,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []Field{
					{Data: []byte("1")},
					{Data: []byte("2.2")},
					{Data: []byte("(5:three)")},
					{Data: []byte("40000")},
				},
			},
			{
				Header: Header{
					Version:       1,
					TotalField:    6,
					Depth:         0,
					Key:           102,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []Field{
					{Data: []byte("")},
					{Data: []byte("(3:abc)")},
					{Data: []byte("(3:foo)")},
					{Data: []byte("(3:bar)")},
					{Data: []byte("5.0")},
					{Data: []byte("")},
					{Data: []byte("777")},
				},
			},
		},
	}

	container0 := mdc.GetContainer(101)
	assert.Equal(t, mdc.Containers[0], *container0)
	assert.Equal(t, "1", container0.GetField(0).String())
	assert.Equal(t, "2.2", container0.GetField(1).String())
	assert.Equal(t, "(5:three)", container0.GetField(2).String())
	assert.Equal(t, "40000", container0.GetField(3).String())

	container1 := mdc.GetContainer(102)
	assert.Equal(t, mdc.Containers[1], *container1)
	assert.Equal(t, "", container1.GetField(0).String())
	assert.Equal(t, "(3:abc)", container1.GetField(1).String())
	assert.Equal(t, "(3:foo)", container1.GetField(2).String())
	assert.Equal(t, "(3:bar)", container1.GetField(3).String())
	assert.Equal(t, "5.0", container1.GetField(4).String())
	assert.Equal(t, "", container1.GetField(5).String())
	assert.Equal(t, "777", container1.GetField(6).String())

	// Field 7 does not exist
	assert.Equal(t, "", container1.GetField(7).String())
	assert.Equal(t, true, container1.GetField(7).IsNull)

	// Container 2 does not exist
	container2 := mdc.GetContainer(1000)
	assert.Nil(t, container2)
}

func TestSetContainer(t *testing.T) {
	containers := Containers{
		Containers: []Container{
			{
				Header: Header{
					Version:       1,
					TotalField:    25,
					Depth:         0,
					Key:           93,
					SchemaVersion: 5263,
					ExtVersion:    13,
				},
				Fields: []Field{
					*NewBasicField("aaa"),
					*NewNullField(field.DateTime),
					*NewBasicField(int32(-1877540863)),
				},
			},
		},
	}
	container0 := containers.GetContainer(93)

	container0.SetField(0, &Field{Data: []byte("(3:bbb)"), Value: "bbb"})
	// t.Logf("containers: \n%s\n", containers.Dump())
	field0, err := containers.GetContainer(93).GetField(0).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, "bbb", field0)

	container0.SetField(18, &Field{Data: []byte(strconv.Itoa(2001)), Value: int32(2001)})
	// t.Logf("containers: \n%s\n", containers.Dump())
	field18, err := containers.GetContainer(93).GetField(18).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, int32(2001), field18)
}

func TestCastVersion(t *testing.T) {

	container := &Container{
		Header: Header{
			Version:       1,
			TotalField:    25,
			Depth:         0,
			Key:           10101,
			SchemaVersion: 5263,
			ExtVersion:    13,
		},
		Fields: []Field{
			*NewBasicField("aaa"),
			*NewNullField(field.DateTime),
			*NewBasicField(int32(-1877540863)),
		},
		Definition: &dictionary.ContainerDefinition{
			Key:           10101,
			Name:          "TestContainer",
			SchemaVersion: 5263,
			ExtVersion:    13,
			Fields: []dictionary.FieldDefinition{
				{Name: "Field1", Type: field.String},
				{Name: "Field2", Type: field.DateTime},
				{Name: "Field3", Type: field.Int32},
			},
		},
	}
	containers := Containers{
		Containers: []Container{
			*container,
		},
	}

	t.Logf("Original containers \n%s\n", containers.Dump())

	definition := &dictionary.ContainerDefinition{
		Key:           10101,
		Name:          "TestContainer",
		SchemaVersion: 5270,
		ExtVersion:    1,
		Fields: []dictionary.FieldDefinition{
			{Name: "Field1", Type: field.String},
			{Name: "NewField1", Type: field.String},
			{Name: "Field2", Type: field.DateTime},
			{Name: "Field3", Type: field.Int32},
		},
	}

	definitions := dictionary.New()
	definitions.Add(definition)

	newContainers, err := containers.CastVersion(definitions, 5270, 1)
	assert.Nil(t, err)
	t.Logf("New version casted containers: \n%s\n", newContainers.Dump())

	assert.Equal(t, 1, len(newContainers.Containers))
	newContainer := newContainers.GetContainer(10101)
	assert.NotNil(t, newContainer)

	assert.Equal(t, 10101, newContainer.Header.Key)
	assert.Equal(t, 5270, newContainer.Header.SchemaVersion)
	assert.Equal(t, 1, newContainer.Header.ExtVersion)
	assert.Equal(t, 4, newContainer.Header.TotalField)
	assert.Equal(t, 4, len(newContainer.Fields))

	v, err := newContainer.GetField(0).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, "aaa", v)

	v, err = newContainer.GetField(1).GetValue()
	assert.Nil(t, err)
	assert.Nil(t, v)

	v, err = newContainer.GetField(2).GetValue()
	assert.Nil(t, err)
	assert.Nil(t, v)

	v, err = newContainer.GetField(3).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, int32(-1877540863), v)
}

func TestCastVersionWithOldVersion(t *testing.T) {

	container := &Container{
		Header: Header{
			Version:       1,
			TotalField:    3,
			Depth:         0,
			Key:           10101,
			SchemaVersion: 5263,
			ExtVersion:    13,
		},
		Fields: []Field{
			*NewBasicField("aaa"),
			*NewNullField(field.DateTime),
			*NewBasicField(int32(-1877540863)),
		},
		Definition: &dictionary.ContainerDefinition{
			Key:           10101,
			Name:          "TestContainer",
			SchemaVersion: 5263,
			ExtVersion:    13,
			Fields: []dictionary.FieldDefinition{
				{Name: "Field1", Type: field.String},
				{Name: "Field2", Type: field.DateTime},
				{Name: "Field3", Type: field.Int32},
			},
		},
	}
	containers := Containers{
		Containers: []Container{
			*container,
		},
	}

	t.Logf("Original containers \n%s\n", containers.Dump())

	// old schema with lesser field
	matrixxSchema := &dictionary.Configuration{
		Name:     xml.Name{Space: "", Local: "configuration"},
		Subtypes: nil,
		Containers: []dictionary.Container{
			{
				ID:                   "TestContainer",
				DocDescription:       "Test Container",
				Key:                  10101,
				CreatedSchemaVersion: 5260,
				DeletedSchemaVersion: 0,
				Fields: []dictionary.Field{
					{
						ID:                   "Field1",
						DocDescription:       "Field 1",
						Datatype:             "string",
						IsArray:              false,
						IsList:               false,
						StructID:             "",
						SubTypeReference:     "",
						CreatedSchemaVersion: 5263,
						DeletedSchemaVersion: 0,
					},
					{
						ID:                   "Field2",
						DocDescription:       "Field 2",
						Datatype:             "datetime",
						IsArray:              false,
						IsList:               false,
						StructID:             "",
						SubTypeReference:     "",
						CreatedSchemaVersion: 5263,
						DeletedSchemaVersion: 0,
					},
				},
			},
		},
	}
	definitions := dictionary.NewWithSchema(matrixxSchema, nil)

	newContainers, err := containers.CastVersion(definitions, 5263, 13)
	assert.Nil(t, err)
	t.Logf("New version casted containers: \n%s\n", newContainers.Dump())

	assert.Equal(t, 1, len(newContainers.Containers))
	newContainer := newContainers.GetContainer(10101)
	assert.NotNil(t, newContainer)

	assert.Equal(t, 10101, newContainer.Header.Key)
	assert.Equal(t, 5263, newContainer.Header.SchemaVersion)
	assert.Equal(t, 13, newContainer.Header.ExtVersion)
	assert.Equal(t, 3, newContainer.Header.TotalField)
	assert.Equal(t, 3, len(newContainer.Fields))

	v, err := newContainer.GetField(0).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, "aaa", v)

	v, err = newContainer.GetField(1).GetValue()
	assert.Nil(t, err)
	assert.Nil(t, v)

	v, err = newContainer.GetField(2).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, int32(-1877540863), v)
}
