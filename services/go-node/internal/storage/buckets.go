package storage

const (
	// bucket names
	bucketBlocks    = "blocks"     // key: blockHash -> value: serialized block bytes
	bucketChainMeta = "chain_meta" // key: "head" -> value: headHash
)
