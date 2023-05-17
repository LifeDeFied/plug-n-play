package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"

    "github.com/cosmos/cosmos-sdk/client"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/gorilla/mux"
)

type Asset struct {
    Owner   string `json:"owner"`
    TokenID string `json:"token_id"`
}

type Donation struct {
    Donor    string `json:"donor"`
    Asset    Asset  `json:"asset"`
    Amount   int64  `json:"amount"`
    Currency string `json:"currency"`
}

var donations []Donation

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/donations", getDonations).Methods("GET")
    router.HandleFunc("/donations", createDonation).Methods("POST")

    http.ListenAndServe(":8000", router)
}

// Call a blockchain RPC via http to retrieve donations.
func getDonations(w http.ResponseWriter, r *http.Request) {
    // Request data from the JSON-RPC server with method "custom/donations" and empty parameters array
    jsonData := map[string]interface{}{
        "jsonrpc": "2.0",
        "id":      "1",
        "method":  "custom/donations",
        "params":  []interface{}{},
    }
    jsonValue, _ := json.Marshal(jsonData)
    response, err := http.Post("http://localhost:26657", "application/json", bytes.NewBuffer(jsonValue))
    if err != nil {
        fmt.Println(err.Error())
    }
    defer response.Body.Close()

    // Decode the response into a slice of donations
    var result interface{}
    body, _ := ioutil.ReadAll(response.Body)
    json.Unmarshal([]byte(body), &result)

    // Write the donations to the response
    json.NewEncoder(w).Encode(result)
}

// Broadcast donations through the network.
func createDonation(w http.ResponseWriter, r *http.Request) {
    reqBody, _ := ioutil.ReadAll(r.Body)
    var donation Donation
    json.Unmarshal(reqBody, &donation)

    // Build the transaction
    clientCtx := client.Context{}.WithCodec(cdc)
    txBuilder := auth.NewTxBuilder(auth.DefaultTxEncoder(), 0, 0)

    fromAddress := "YOUR_SENDING_ADDRESS"
    toAddress := "YOUR_RECEIVING_ADDRESS"

    msg := &types.MsgSend{
        FromAddress: sdk.AccAddress(fromAddress),
        ToAddress:   sdk.AccAddress(toAddress),
        Amount:      sdk.NewCoins(sdk.NewCoin(donation.Currency, sdk.NewInt(donation.Amount))),
    }

    err := txBuilder.SetMsgs(msg)
    if err != nil {
        rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
        return
    }

    // Sign the transaction
    err = txBuilder.Sign("donor", keyManager.GetAddr().String(), auth.StdSignMsg{})
    if err != nil {
        rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
        return
    }

    // Get the signed transaction
    txBytes, err := clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
    if err != nil {
        rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
        return
    }

    // Broadcast the transaction
    res, err := clientCtx.BroadcastTx(txBytes)
    if err != nil {
        rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
        return
    }

    donations = append(donations, donation)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(res)
}
