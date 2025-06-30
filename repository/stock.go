package repository

type itemStock []OrderItem

func LoadItemStock() itemStock {
	return []OrderItem{
		{Id: "1", Description: "Whey", Price: 12490},
		{Id: "2", Description: "Creatina", Price: 6490},
		{Id: "3", Description: "Multivitamínico", Price: 4790},
		{Id: "4", Description: "Pré Treino", Price: 10490},
		{Id: "5", Description: "Beta Alanina", Price: 6890},
		{Id: "6", Description: "Barrinha", Price: 2990},
		{Id: "7", Description: "Glutamina", Price: 4990},
		{Id: "8", Description: "Cafeína", Price: 4490},
		{Id: "9", Description: "Melatonina", Price: 3390},
		{Id: "10", Description: "Óleo de Peixe", Price: 5290},
	}
}

func (i itemStock) Get(id int) OrderItem {
	return i[id-1]
}
