package domain

type Ticket struct {
	ID      string `bson:"_id,omitempty"`
	Title   string `bson:"title"`
	Content string `bson:"content"`
	Pie     string `bson:"pie"`
	Numero  int    `bson:"numero"`
	LogoURL string `bson:"logo_url"`
}

type Texto struct {
	ID    string `bson:"_id,omitempty"`
	Clave string `bson:"clave"`
	Valor string `bson:"valor"`
}

type Logo struct {
	ID      string `bson:"_id,omitempty"`
	FileURL string `bson:"file_url"`
}
