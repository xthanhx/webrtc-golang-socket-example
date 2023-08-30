package event

import (
	"encoding/json"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Info struct {
	PosId     string    `json:"pos_id"`
	Role      string    `json:"role"`
	Name      string    `json:"name"`
	ConnectID uuid.UUID `json:"connect_id"`
}

type InfoReq struct {
	PosId string `json:"pos_id"`
	Role  string `json:"role"`
	Name  string `json:"name"`
}

type requestSignal struct {
	ConnectID uuid.UUID `json:"connect_id"`
}

type verifyConnect struct {
	ConnectID    uuid.UUID `json:"connect_id"`
	NumberVerify int       `json:"number_verify"`
}
type signalConnect struct {
	RoomChannel string `json:"room_channel"`
	payload     any    `json:"payload"`
}

type roomChannelConnect struct {
	RoomChannel string `json:"room_channel"`
}

var signalPath = "/signal"

func generateRoomName(posId string) string {
	return "pos_channel_" + posId
}

func randomNumber() int {
	rand.Seed(time.Now().UnixNano())
	max := 100
	return rand.Intn(max)
}

func includesElement(arr []int, target int) bool {
	for _, val := range arr {
		if val == target {
			return true
		}
	}
	return false
}

func randomNumbers(amount int) []int {
	numbers := make([]int, 0)

	for len(numbers) < amount {
		rdn := randomNumber()
		if !includesElement(numbers, rdn) {
			numbers = append(numbers, rdn)
		}
	}

	return numbers
}

func shuffleArray(array []int) []int {
	shuffled := make([]int, len(array))
	copy(shuffled, array)

	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled
}

func mapData(payload any, data any) error {
	jsonStr, _ := json.Marshal(payload)
	err := json.Unmarshal(jsonStr, &data)
	if err != nil {
		return err
	}
	return nil
}
