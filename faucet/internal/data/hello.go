package data

type HelloResponse struct {
	Epoch      uint64 `json:"epoch"`
	PeerNumber int    `json:"peer_number"`
}
