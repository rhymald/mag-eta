package functions 

import (
	"fmt"
	"crypto/sha512"
	"encoding/binary"
	"github.com/btcsuite/btcutil/base58"
)

func GetID(base, sess int) string { 
	in_bytes := make([]byte, 8)
	// bid := (*c).TSBorn
	// aid := (*c).TSAtts
	binary.LittleEndian.PutUint64(in_bytes, uint64(base))
	str := sha512.Sum512(in_bytes)
	bornID := base58.Encode(str[:])
	binary.LittleEndian.PutUint64(in_bytes, uint64(sess))
	str = sha512.Sum512(in_bytes)
	attsID := base58.Encode(str[:])
	return fmt.Sprintf("%v-%v-%v-%v", bornID[:4], bornID[(len(bornID)-9):len(bornID)], attsID[:1], attsID[(len(attsID)-9):len(attsID)])
	// return bornID
	// return (*c).ID["Born"] 
}
