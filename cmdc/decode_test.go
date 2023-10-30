package cmdc

import (
	"errors"
	"testing"

	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/stretchr/testify/assert"
)

func TestDecodeSingleContainer1(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,20,300,4]"
	containers, err := Decode([]byte(mdc))
	assert.Nil(t, err)

	container := containers.Containers[0]

	expectedHeader := mdd.Header{
		Version:       1,
		TotalField:    18,
		Depth:         0,
		Key:           -6,
		SchemaVersion: 5222,
		ExtVersion:    2,
	}
	assert.Equal(t, expectedHeader, container.Header)
	assert.Equal(t, "1", container.GetField(0).String())
	assert.Equal(t, "20", container.GetField(1).String())
	assert.Equal(t, "300", container.GetField(2).String())
	assert.Equal(t, "4", container.GetField(3).String())
}

func TestDecodeSingleContainer2(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[,(6:value2),3,2021-09-07T08:00:25.000001Z," +
		"2021-10-31,09:13:02.667997Z,88,5.5,]"
	containers, err := Decode([]byte(mdc))
	assert.Nil(t, err)

	container := containers.Containers[0]

	expectedHeader := mdd.Header{
		Version:       1,
		TotalField:    18,
		Depth:         0,
		Key:           -6,
		SchemaVersion: 5222,
		ExtVersion:    2,
	}
	assert.Equal(t, expectedHeader, container.Header)
	assert.Equal(t, "", container.GetField(0).String())
	assert.Equal(t, "(6:value2)", container.GetField(1).String())
	assert.Equal(t, "3", container.GetField(2).String())
	assert.Equal(t, "2021-09-07T08:00:25.000001Z", container.GetField(3).String())
	assert.Equal(t, "2021-10-31", container.GetField(4).String())
	assert.Equal(t, "09:13:02.667997Z", container.GetField(5).String())
	assert.Equal(t, "88", container.GetField(6).String())
	assert.Equal(t, "5.5", container.GetField(7).String())
	assert.Equal(t, "", container.GetField(8).String())
}

func TestDecodeContainers(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,20,300,4]<1,5,0,-7,5222,2>[,2,(3:def),4]"
	containers, err := Decode([]byte(mdc))
	assert.Nil(t, err)

	container0 := containers.Containers[0]

	expectedHeader := mdd.Header{
		Version:       1,
		TotalField:    18,
		Depth:         0,
		Key:           -6,
		SchemaVersion: 5222,
		ExtVersion:    2,
	}
	assert.Equal(t, expectedHeader, container0.Header)
	assert.Equal(t, "1", container0.GetField(0).String())
	assert.Equal(t, "20", container0.GetField(1).String())
	assert.Equal(t, "300", container0.GetField(2).String())
	assert.Equal(t, "4", container0.GetField(3).String())

	container1 := containers.Containers[1]

	expectedHeader = mdd.Header{
		Version:       1,
		TotalField:    5,
		Depth:         0,
		Key:           -7,
		SchemaVersion: 5222,
		ExtVersion:    2,
	}
	assert.Equal(t, expectedHeader, container1.Header)
	assert.Equal(t, "", container1.GetField(0).String())
	assert.Equal(t, "2", container1.GetField(1).String())
	assert.Equal(t, "(3:def)", container1.GetField(2).String())
	assert.Equal(t, "4", container1.GetField(3).String())
}

func TestDecodeNestedContainers(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,20,<1,2,0,452,5222,2>[100],4]"
	containers, err := Decode([]byte(mdc))
	assert.Nil(t, err)

	container0 := containers.Containers[0]
	assert.Equal(t, "1", container0.GetField(0).String())
	assert.Equal(t, "20", container0.GetField(1).String())
	assert.Equal(t, "<1,2,0,452,5222,2>[100]", container0.GetField(2).String())
	assert.Equal(t, "4", container0.GetField(3).String())
}

func TestDecodeFieldWithReservedCharacter(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,2,(10:v[<ue(obar),4,,6]"
	containers, err := Decode([]byte(mdc))
	assert.Nil(t, err)

	container0 := containers.Containers[0]
	assert.Equal(t, "1", container0.GetField(0).String())
	assert.Equal(t, "2", container0.GetField(1).String())
	assert.Equal(t, "(10:v[<ue(obar)", container0.GetField(2).String())
	assert.Equal(t, "4", container0.GetField(3).String())
	assert.Equal(t, "", container0.GetField(4).String())
	assert.Equal(t, "6", container0.GetField(5).String())
}

func TestDecodeNestedContainersWithReservedCharacter(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,2,<1,2,0,452,5222,2>[100,(5:a[<,(),300],3]"
	containers, err := Decode([]byte(mdc))
	assert.Nil(t, err)

	container0 := containers.Containers[0]
	assert.Equal(t, "1", container0.GetField(0).String())
	assert.Equal(t, "2", container0.GetField(1).String())
	assert.Equal(t, "<1,2,0,452,5222,2>[100,(5:a[<,(),300]", container0.GetField(2).String())
	assert.Equal(t, "3", container0.GetField(3).String())
}

func TestDecodeEmptyBody(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[]"
	containers, err := Decode([]byte(mdc))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(containers.Containers))
	container0 := containers.Containers[0]
	assert.Equal(t, "", container0.GetField(0).String())
}

func TestInvalidHeader(t *testing.T) {
	mdc := "<1,18,-6,5222,2>[1,20,300,4]"
	_, err := Decode([]byte(mdc))
	assert.Equal(t, errors.New("Invalid cMDC header, 6 fields expected"), err)
}

func TestInvalidHeader2(t *testing.T) {
	mdc := "<1,18,0,-6,5222[1,20,300,4]"
	_, err := Decode([]byte(mdc))
	assert.Equal(t, errors.New("Invalid cMDC character '[' in header, numeric expected"), err)
}

func TestInvalidHeader3(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2"
	_, err := Decode([]byte(mdc))
	assert.Equal(t, errors.New("Invalid cMDC header, no end of header"), err)
}

func TestInvalidHeader4(t *testing.T) {
	mdc := "1,18,0,-6,5222,2>[]"
	_, err := Decode([]byte(mdc))
	assert.Equal(t, errors.New("Invalid cMDC header, first character must be '<'"), err)
}

func TestInvalidHeader5(t *testing.T) {
	mdc := "<1,18,0,1-6,5222,2>[]"
	_, err := Decode([]byte(mdc))
	assert.Equal(t, errors.New("Invalid cMDC header field '1-6', numeric expected"), err)
}

func TestInvalidBody(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,20,300,4"
	_, err := Decode([]byte(mdc))
	assert.Equal(t, errors.New("Invalid cMDC body, no end of body"), err)
}

func TestInvalidBody2(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>"
	_, err := Decode([]byte(mdc))
	assert.Equal(t, errors.New("Invalid cMDC body, no body"), err)
}

func TestInvalidBody3(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>1,2,3]"
	_, err := Decode([]byte(mdc))
	assert.Equal(t, errors.New("Invalid cMDC body, first character must be '['"), err)
}

func TestInvalidBody4(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,(abc:foo),3,4]"
	_, err := Decode([]byte(mdc))
	assert.Equal(t, errors.New("Invalid character 'a', numeric expected for string length"), err)
}