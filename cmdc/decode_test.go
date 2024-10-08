package cmdc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matrixxsoftware/go-mdd/mdd"
)

func TestDecodeSingleContainer1(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,20,300,4]"
	containers, err := codec.Decode([]byte(mdc))
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
	containers, err := codec.Decode([]byte(mdc))
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
	containers, err := codec.Decode([]byte(mdc))
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
	mdc := "<1,18,0,-6,5222,2>[1,20,<1,2,0,452,5222,2>[100,<1,2,0,450,5222,2>[999]],4]"
	containers, err := codec.Decode([]byte(mdc))
	assert.Nil(t, err)

	container0 := containers.Containers[0]
	assert.Equal(t, "1", container0.GetField(0).String())
	assert.Equal(t, "20", container0.GetField(1).String())
	assert.Equal(t, "<1,2,0,452,5222,2>[100,<1,2,0,450,5222,2>[999]]", container0.GetField(2).String())
	assert.Equal(t, "4", container0.GetField(3).String())

	assert.Equal(t, false, container0.GetField(0).IsContainer)
	assert.Equal(t, false, container0.GetField(1).IsContainer)
	assert.Equal(t, true, container0.GetField(2).IsContainer)
	assert.Equal(t, false, container0.GetField(3).IsContainer)
}

func TestDecodeMultipleNestedContainers(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,20,<1,2,0,452,5222,2>[100],4,<1,2,0,451,5222,2>[666]]"
	containers, err := codec.Decode([]byte(mdc))
	assert.Nil(t, err)

	container0 := containers.Containers[0]
	assert.Equal(t, "1", container0.GetField(0).String())
	assert.Equal(t, "20", container0.GetField(1).String())
	assert.Equal(t, "<1,2,0,452,5222,2>[100]", container0.GetField(2).String())
	assert.Equal(t, "4", container0.GetField(3).String())
	assert.Equal(t, "<1,2,0,451,5222,2>[666]", container0.GetField(4).String())

	assert.Equal(t, false, container0.GetField(0).IsContainer)
	assert.Equal(t, false, container0.GetField(1).IsContainer)
	assert.Equal(t, true, container0.GetField(2).IsContainer)
	assert.Equal(t, false, container0.GetField(3).IsContainer)
	assert.Equal(t, true, container0.GetField(4).IsContainer)
}

func TestDecodeNestedContainersWithList(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,20,<1,2,0,452,5222,2>[100,{1,2,3}],4]"
	containers, err := codec.Decode([]byte(mdc))
	assert.Nil(t, err)

	container0 := containers.Containers[0]
	assert.Equal(t, "1", container0.GetField(0).String())
	assert.Equal(t, "20", container0.GetField(1).String())
	assert.Equal(t, "<1,2,0,452,5222,2>[100,{1,2,3}]", container0.GetField(2).String())
	assert.Equal(t, "4", container0.GetField(3).String())

	assert.Equal(t, false, container0.GetField(0).IsContainer)
	assert.Equal(t, false, container0.GetField(1).IsContainer)
	assert.Equal(t, true, container0.GetField(2).IsContainer)
	assert.Equal(t, false, container0.GetField(3).IsContainer)

	assert.Equal(t, false, container0.GetField(0).IsMulti)
	assert.Equal(t, false, container0.GetField(1).IsMulti)
	assert.Equal(t, false, container0.GetField(2).IsMulti)
	assert.Equal(t, false, container0.GetField(3).IsMulti)
}

func TestDecodeNestedContainersWithReservedCharacter(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,2,<1,2,0,452,5222,2>[100,(5:a[<,(),300],3]"
	containers, err := codec.Decode([]byte(mdc))
	assert.Nil(t, err)

	container0 := containers.Containers[0]
	assert.Equal(t, "1", container0.GetField(0).String())
	assert.Equal(t, "2", container0.GetField(1).String())
	assert.Equal(t, "<1,2,0,452,5222,2>[100,(5:a[<,(),300]", container0.GetField(2).String())
	assert.Equal(t, "3", container0.GetField(3).String())
}

func TestIsNull(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[,(6:foobar),,4,]"
	containers, err := codec.Decode([]byte(mdc))
	assert.Nil(t, err)

	container := containers.Containers[0]

	assert.Equal(t, 5, len(container.Fields))
	assert.Equal(t, "", container.GetField(0).String())
	assert.Equal(t, "(6:foobar)", container.GetField(1).String())
	assert.Equal(t, "", container.GetField(2).String())
	assert.Equal(t, "4", container.GetField(3).String())
	assert.Equal(t, "", container.GetField(4).String())

	assert.Equal(t, true, container.GetField(0).IsNull)
	assert.Equal(t, false, container.GetField(1).IsNull)
	assert.Equal(t, true, container.GetField(2).IsNull)
	assert.Equal(t, false, container.GetField(3).IsNull)
	assert.Equal(t, true, container.GetField(4).IsNull)
}

func TestListIntegerValue(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[0,{1,2,3},,,300,{4,5}]"
	containers, err := codec.Decode([]byte(mdc))
	assert.Nil(t, err)

	container0 := containers.Containers[0]
	assert.Equal(t, "0", container0.GetField(0).String())
	assert.Equal(t, "{1,2,3}", container0.GetField(1).String())
	assert.Equal(t, "", container0.GetField(2).String())
	assert.Equal(t, "", container0.GetField(3).String())
	assert.Equal(t, "300", container0.GetField(4).String())
	assert.Equal(t, "{4,5}", container0.GetField(5).String())

	assert.Equal(t, false, container0.GetField(0).IsMulti)
	assert.Equal(t, true, container0.GetField(1).IsMulti)
	assert.Equal(t, false, container0.GetField(2).IsMulti)
	assert.Equal(t, false, container0.GetField(3).IsMulti)
	assert.Equal(t, false, container0.GetField(4).IsMulti)
	assert.Equal(t, true, container0.GetField(5).IsMulti)
}

func TestListStringValue(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[0,{(3:one),(3:two),(5:three)},,,300,{(4:four),(4:five)}]"
	containers, err := codec.Decode([]byte(mdc))
	assert.Nil(t, err)

	container0 := containers.Containers[0]
	assert.Equal(t, "0", container0.GetField(0).String())
	assert.Equal(t, "{(3:one),(3:two),(5:three)}", container0.GetField(1).String())
	assert.Equal(t, "", container0.GetField(2).String())
	assert.Equal(t, "", container0.GetField(3).String())
	assert.Equal(t, "300", container0.GetField(4).String())
	assert.Equal(t, "{(4:four),(4:five)}", container0.GetField(5).String())

	assert.Equal(t, false, container0.GetField(0).IsMulti)
	assert.Equal(t, true, container0.GetField(1).IsMulti)
	assert.Equal(t, false, container0.GetField(2).IsMulti)
	assert.Equal(t, false, container0.GetField(3).IsMulti)
	assert.Equal(t, false, container0.GetField(4).IsMulti)
	assert.Equal(t, true, container0.GetField(5).IsMulti)
}

func TestListStructValue(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[0,{<1,3,0,100,5222,2>[,1000,,],<1,3,0,100,5222,2>[,2000,,],<1,3,0,100,5222,2>[,3000,]},,,{300},{<1,5,0,10,5222,2>[4000]}]"
	containers, err := codec.Decode([]byte(mdc))
	assert.Nil(t, err)
	// t.Logf(containers.Dump())

	container0 := containers.Containers[0]
	assert.Equal(t, "0", container0.GetField(0).String())
	assert.Equal(t, "{<1,3,0,100,5222,2>[,1000,,],<1,3,0,100,5222,2>[,2000,,],<1,3,0,100,5222,2>[,3000,]}", container0.GetField(1).String())
	assert.Equal(t, "", container0.GetField(2).String())
	assert.Equal(t, "", container0.GetField(3).String())
	assert.Equal(t, "{300}", container0.GetField(4).String())
	assert.Equal(t, "{<1,5,0,10,5222,2>[4000]}", container0.GetField(5).String())

	assert.Equal(t, false, container0.GetField(0).IsMulti)
	assert.Equal(t, true, container0.GetField(1).IsMulti)
	assert.Equal(t, false, container0.GetField(2).IsMulti)
	assert.Equal(t, false, container0.GetField(3).IsMulti)
	assert.Equal(t, true, container0.GetField(4).IsMulti)
	assert.Equal(t, true, container0.GetField(5).IsMulti)

	assert.Equal(t, false, container0.GetField(0).IsContainer)
	assert.Equal(t, true, container0.GetField(1).IsContainer)
	assert.Equal(t, false, container0.GetField(2).IsContainer)
	assert.Equal(t, false, container0.GetField(3).IsContainer)
	assert.Equal(t, false, container0.GetField(4).IsContainer)
	assert.Equal(t, true, container0.GetField(5).IsContainer)
}

func TestDecodeEmptyBody(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[]"
	containers, err := codec.Decode([]byte(mdc))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(containers.Containers))
	container0 := containers.Containers[0]
	assert.Equal(t, "", container0.GetField(0).String())
}

func TestZeroLenStringField(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,(0:),3,4]"
	containers, err := codec.Decode([]byte(mdc))
	assert.Nil(t, err)

	container := containers.Containers[0]
	assert.Equal(t, "1", container.GetField(0).String())
	assert.Equal(t, "(0:)", container.GetField(1).String())
	assert.Equal(t, "3", container.GetField(2).String())
	assert.Equal(t, "4", container.GetField(3).String())
}

func TestEmptyStringField(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,(),3,4]"
	containers, err := codec.Decode([]byte(mdc))
	assert.Nil(t, err)

	container := containers.Containers[0]
	assert.Equal(t, "1", container.GetField(0).String())
	assert.Equal(t, "()", container.GetField(1).String())
	assert.Equal(t, "3", container.GetField(2).String())
	assert.Equal(t, "4", container.GetField(3).String())
}

func TestUnicodeStringField(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,(6:富爸),3,4]"
	containers, err := codec.Decode([]byte(mdc))
	assert.Nil(t, err)

	container := containers.Containers[0]
	assert.Equal(t, "1", container.GetField(0).String())
	assert.Equal(t, "(6:富爸)", container.GetField(1).String())
	assert.Equal(t, "3", container.GetField(2).String())
	assert.Equal(t, "4", container.GetField(3).String())
}

func TestEscapeStringField(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,(6:ab\\\\def),3,4]"
	containers, err := codec.Decode([]byte(mdc))
	assert.Nil(t, err)

	container := containers.Containers[0]
	assert.Equal(t, "1", container.GetField(0).String())
	assert.Equal(t, "(6:ab\\\\def)", container.GetField(1).String())
	assert.Equal(t, "3", container.GetField(2).String())
	assert.Equal(t, "4", container.GetField(3).String())
}

func TestBlobField(t *testing.T) {
	mdc := "<1,18,0,-999,5222,2>[1,(17:\\c2\\82\"\\c3\\b8\\10\\08I\"\\c3\\b8\\10\\01\\01\\c3\\85\\05),3,4]"
	containers, err := codec.Decode([]byte(mdc))
	assert.Nil(t, err)

	container := containers.Containers[0]
	assert.Equal(t, "1", container.GetField(0).String())
	assert.Equal(t, "(17:\\c2\\82\"\\c3\\b8\\10\\08I\"\\c3\\b8\\10\\01\\01\\c3\\85\\05)", container.GetField(1).String())
	assert.Equal(t, "3", container.GetField(2).String())
	assert.Equal(t, "4", container.GetField(3).String())
}

func TestReservedCharacterStringField(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,2,(10:v[<ue(obar),4,,6]"
	containers, err := codec.Decode([]byte(mdc))
	assert.Nil(t, err)

	container0 := containers.Containers[0]
	assert.Equal(t, "1", container0.GetField(0).String())
	assert.Equal(t, "2", container0.GetField(1).String())
	assert.Equal(t, "(10:v[<ue(obar)", container0.GetField(2).String())
	assert.Equal(t, "4", container0.GetField(3).String())
	assert.Equal(t, "", container0.GetField(4).String())
	assert.Equal(t, "6", container0.GetField(5).String())
}

func TestInvalidHeader(t *testing.T) {
	mdc := "<1,18,-6,5222,2>[1,20,300,4]"
	_, err := codec.Decode([]byte(mdc))
	assert.Equal(t, errors.New("invalid cMDC header, 6 fields expected"), err)
}

func TestInvalidHeader2(t *testing.T) {
	mdc := "<1,18,0,-6,5222[1,20,300,4]"
	_, err := codec.Decode([]byte(mdc))
	assert.Equal(t, errors.New("invalid cMDC character '[' in header, numeric expected"), err)
}

func TestInvalidHeader3(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2"
	_, err := codec.Decode([]byte(mdc))
	assert.Equal(t, errors.New("invalid cMDC header, no end of header"), err)
}

func TestInvalidHeader4(t *testing.T) {
	mdc := "1,18,0,-6,5222,2>[]"
	_, err := codec.Decode([]byte(mdc))
	assert.Equal(t, errors.New("invalid cMDC header, first character must be '<'"), err)
}

func TestInvalidHeader5(t *testing.T) {
	mdc := "<1,18,0,1-6,5222,2>[]"
	_, err := codec.Decode([]byte(mdc))
	assert.Equal(t, errors.New("invalid cMDC header field '1-6', numeric expected"), err)
}

func TestInvalidBody(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,20,300,4"
	_, err := codec.Decode([]byte(mdc))
	assert.Equal(t, errors.New("invalid cMDC body, no end of body"), err)
}

func TestInvalidBody2(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>"
	_, err := codec.Decode([]byte(mdc))
	assert.Equal(t, errors.New("invalid cMDC body, no body"), err)
}

func TestInvalidBody3(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>1,2,3]"
	_, err := codec.Decode([]byte(mdc))
	assert.Equal(t, errors.New("invalid cMDC body, first character must be '['"), err)
}

func TestInvalidBody4(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,(abc:foo),3,4]"
	_, err := codec.Decode([]byte(mdc))
	assert.Equal(t, errors.New("invalid character 'a', numeric expected for string length"), err)
}

func TestInvalidBody5(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,(5:foo),3,4]"
	_, err := codec.Decode([]byte(mdc))
	assert.Equal(t, errors.New("invalid cMDC body, mismatch string length"), err)
}

func TestInvalidBody6(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,(5:foobar),3,4]"
	_, err := codec.Decode([]byte(mdc))
	assert.Equal(t, errors.New("invalid cMDC body, mismatch string length"), err)
}

func TestInvalidBody7(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[1,(5:fooba:),3,4]"
	_, err := codec.Decode([]byte(mdc))
	assert.Equal(t, errors.New("invalid cMDC body, mismatch string length"), err)
}

func TestInvalidBody8(t *testing.T) {
	mdc := "<1,18,0,-999,5222,2>[1,(31:\\c2\\82\"\\c3\\b8\\10\\08I\"\\c3\\b8\\10\\01\\01\\c3\\85\\05),3,4]"
	_, err := codec.Decode([]byte(mdc))
	assert.Equal(t, errors.New("invalid cMDC body, no end of body"), err)
}

func TestInvalidBody9(t *testing.T) {
	mdc := "<1,18,0,-999,5222,2>[1,(10:\\c2\\82\"\\c3\\b8\\10\\08I\"\\c3\\b8\\10\\01\\01\\c3\\85\\05,3,4]"
	_, err := codec.Decode([]byte(mdc))
	assert.Equal(t, errors.New("invalid cMDC body, mismatch string length"), err)
}

func TestInvalidBody10(t *testing.T) {
	mdc := "<1,18,0,-999,5222,2>[1,(7:富爸),3,4"
	_, err := codec.Decode([]byte(mdc))
	assert.Equal(t, errors.New("invalid cMDC body, mismatch string length"), err)
}

func TestInvalidBody11(t *testing.T) {
	mdc := "<1,18,0,-999,5222,2>[1,(2:富爸),3,4"
	_, err := codec.Decode([]byte(mdc))
	assert.Equal(t, errors.New("invalid cMDC body, mismatch string length"), err)
}

func TestEmptyBody(t *testing.T) {
	mdc := "<1,18,0,-6,5222,2>[]"
	containers, err := codec.Decode([]byte(mdc))
	assert.Nil(t, err)

	container0 := containers.Containers[0]
	assert.Equal(t, 0, len(container0.Fields))
}
