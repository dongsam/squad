package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto"

	"github.com/crescent-network/crescent/x/liquidity/types"
)

type keysTestSuite struct {
	suite.Suite
}

func TestKeysTestSuite(t *testing.T) {
	suite.Run(t, new(keysTestSuite))
}

func (s *keysTestSuite) TestGetPairKey() {
	s.Require().Equal([]byte{0xa5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, types.GetPairKey(0))
	s.Require().Equal([]byte{0xa5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x9}, types.GetPairKey(9))
	s.Require().Equal([]byte{0xa5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa}, types.GetPairKey(10))
}

func (s *keysTestSuite) TestGetPoolKey() {
	s.Require().Equal([]byte{0xab, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}, types.GetPoolKey(1))
	s.Require().Equal([]byte{0xab, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5}, types.GetPoolKey(5))
	s.Require().Equal([]byte{0xab, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa}, types.GetPoolKey(10))
}

func (s *keysTestSuite) TestPairLookupIndexKey() {
	testCases := []struct {
		denomA   string
		denomB   string
		pairId   uint64
		expected []byte
	}{
		{
			"denomA",
			"denomB",
			uint64(1),
			[]byte{0xa7, 0x6, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x41, 0x6, 0x64,
				0x65, 0x6e, 0x6f, 0x6d, 0x42, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1},
		},
		{
			"denomC",
			"denomD",
			uint64(20),
			[]byte{0xa7, 0x6, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x43, 0x6, 0x64,
				0x65, 0x6e, 0x6f, 0x6d, 0x44, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x14},
		},
		{
			"denomE",
			"denomF",
			uint64(13),
			[]byte{0xa7, 0x6, 0x64, 0x65, 0x6e, 0x6f, 0x6d, 0x45, 0x6, 0x64,
				0x65, 0x6e, 0x6f, 0x6d, 0x46, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xd},
		},
	}

	for _, tc := range testCases {
		key := types.GetPairsByDenomsIndexKey(tc.denomA, tc.denomB, tc.pairId)
		s.Require().Equal(tc.expected, key)

		denomA, denomB, pairId := types.ParsePairsByDenomsIndexKey(key)
		s.Require().Equal(tc.denomA, denomA)
		s.Require().Equal(tc.denomB, denomB)
		s.Require().Equal(tc.pairId, pairId)
	}
}

func (s *keysTestSuite) TestPoolByReserveAccKey() {
	reserveAcc1 := sdk.AccAddress(crypto.AddressHash([]byte("ReserveAccount1")))
	reserveAcc2 := sdk.AccAddress(crypto.AddressHash([]byte("ReserveAccount2")))
	reserveAcc3 := sdk.AccAddress(crypto.AddressHash([]byte("ReserveAccount3")))
	s.Require().Equal([]byte{0xac, 0x14, 0x7c, 0x17, 0x3a, 0x5a, 0xc0, 0xdb, 0x24, 0x94, 0xf9,
		0xd9, 0x81, 0xc6, 0xb6, 0xfe, 0x1d, 0x98, 0x76, 0x47, 0x12, 0x40}, types.GetPoolByReserveAccIndexKey(reserveAcc1))
	s.Require().Equal([]byte{0xac, 0x14, 0xc4, 0xb1, 0xef, 0xd1, 0x16, 0x23, 0xb6, 0x4a, 0x51,
		0xb5, 0xb0, 0x8a, 0xc0, 0xdd, 0x0, 0x71, 0x3f, 0xe3, 0x1f, 0x1d}, types.GetPoolByReserveAccIndexKey(reserveAcc2))
	s.Require().Equal([]byte{0xac, 0x14, 0x7, 0x84, 0x6d, 0x48, 0x46, 0xb2, 0x29, 0x34, 0x3b,
		0x49, 0xc2, 0xd4, 0xee, 0xb5, 0x4d, 0x2, 0x84, 0x50, 0x74, 0x83}, types.GetPoolByReserveAccIndexKey(reserveAcc3))
}

func (s *keysTestSuite) TestPoolsByPairIndexKey() {
	testCases := []struct {
		pairId   uint64
		poolId   uint64
		expected []byte
	}{
		{
			uint64(5),
			uint64(10),
			[]byte{0xad, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa},
		},
		{
			uint64(2),
			uint64(7),
			[]byte{0xad, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7},
		},
		{
			uint64(3),
			uint64(5),
			[]byte{0xad, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5},
		},
	}

	for _, tc := range testCases {
		key := types.GetPoolsByPairIndexKey(tc.pairId, tc.poolId)
		s.Require().Equal(tc.expected, key)

		poolId := types.ParsePoolsByPairIndexKey(key)
		s.Require().Equal(tc.poolId, poolId)
	}
}

func (s *keysTestSuite) TestDepositRequestKey() {
	// TODO: not implemented yet
}

func (s *keysTestSuite) TestWithdrawRequestKey() {
	// TODO: not implemented yet
}

func (s *keysTestSuite) TestSwapRequestKey() {
	// TODO: not implemented yet
}
