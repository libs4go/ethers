package binding

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/libs4go/ethers/abi"
	"github.com/stretchr/testify/require"
)

func TestArrayRegex(t *testing.T) {
	allMatch := ArrayRegex.FindAllString("uint256[1][][3]", -1)

	require.Equal(t, allMatch, []string{"uint256", "[1]", "[]", "[3]"})

	allMatch = ArrayLenRegex.FindStringSubmatch("[]")

	println(fmt.Sprintf("%v", allMatch))
}

func TestTupleNameRegex(t *testing.T) {
	println(fmt.Sprintf(TupleNameRegex.FindStringSubmatch("struct CurveNFT[2][]")[1]))
}

func TestLoadFiles(t *testing.T) {
	symbols := NewSymbols()
	files, err := ioutil.ReadDir("./testdata")

	require.NoError(t, err)

	for _, file := range files {
		_, err = ParseFile(file.Name(), filepath.Join("./testdata", file.Name()), symbols)

		require.NoError(t, err)
	}
}

func TestFooContract(t *testing.T) {
	generator := NewGen()

	contract, err := ParseFile("Foo", "./testdata/Foo.json", generator)

	require.NoError(t, err)

	f, ok := abi.TryGetFunc(contract, "baz(uint32,bool)")

	require.True(t, ok)

	buff, err := f.Call(69, true)

	require.NoError(t, err)

	require.Equal(t, hex.EncodeToString(buff), "cdcd77c000000000000000000000000000000000000000000000000000000000000000450000000000000000000000000000000000000000000000000000000000000001")

	f, ok = abi.TryGetFunc(contract, "bar(bytes3[2])")

	require.True(t, ok)

	var abc [3]byte
	var def [3]byte

	copy(abc[:], []byte("abc"))
	copy(def[:], []byte("def"))

	buff, err = f.Call([2][3]byte{abc, def})

	require.NoError(t, err)

	require.Equal(t, hex.EncodeToString(buff), "fce353f661626300000000000000000000000000000000000000000000000000000000006465660000000000000000000000000000000000000000000000000000000000")

	f, ok = abi.TryGetFunc(contract, "sam(bytes,bool,uint256[])")

	require.True(t, ok)

	buff, err = f.Call([]byte("dave"), true, []uint{1, 2, 3})

	require.NoError(t, err)

	require.Equal(t, hex.EncodeToString(buff), "a5643bf20000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000464617665000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003")

}

func TestGen(t *testing.T) {

	generator := NewGen()

	_, err := ParseFile("Foo", "./testdata/CurveUSDVault.json", generator)

	require.NoError(t, err)

	var writerBuffer bytes.Buffer

	require.NoError(t, generator.Write("testdata", &writerBuffer))

	// println(writerBuffer.String())

	ioutil.WriteFile("./testdata/test.go", writerBuffer.Bytes(), 077)
}

func TestToUpper(t *testing.T) {
	println(strings.Title("hello world"))
}
