package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Cart struct {
	User_ID    int `json:"user_id"`
	Product_ID int `json:"product_id"`
	Quantity   int `json:"quantity"`
}

func addCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}
	username := r.Context().Value(usernameKey)
	if username == nil {
		http.Error(w, "no username in context", http.StatusInternalServerError)
		return
	}
	var c []Cart
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Invalid json format", http.StatusNonAuthoritativeInfo)
		return
	}
	for _, cart := range c {
		row := db.QueryRow("select stock from products where id = $1", cart.Product_ID)
		var stock int
		err := row.Scan(&stock)
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid product id", http.StatusInternalServerError)
			return
		} else if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Panic(err)
			return
		}
		if (stock - cart.Quantity) < 0 {
			http.Error(w, "too much quantity", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec("insert into cart (user_id,product_id,quantity) values ($1,$2,$3)", cart.User_ID, cart.Product_ID, cart.Quantity)
		if err != nil {
			http.Error(w, "failed to insert cart", http.StatusInternalServerError)
			log.Panic(err)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Cart added successfully!",
	})
}

func getCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}
	username := r.Context().Value(usernameKey)
	if username == nil {
		http.Error(w, "no username in context", http.StatusInternalServerError)
		return
	}
	idrow := db.QueryRow("select id from users where username = $1", username)
	var id int
	err := idrow.Scan(&id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Panic(err)
		return
	}
	cartrow := db.QueryRow("select product_id, quantity from cart where user_id = $1", id)
	var c Cart
	err = cartrow.Scan(&c.Product_ID, &c.Quantity)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Panic(err)
		return
	}
	err = json.NewEncoder(w).Encode(c)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Panic(err)
		return
	}
}
