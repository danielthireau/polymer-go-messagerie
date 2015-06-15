package messenger

import (
	"fmt"
	//"appengine"
	"appengine/datastore"
	//"encoding/json"
	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	//"log"
	//"net/http"
)

type MessageAPI struct{}

type Message struct {
	// le json:"name"  datastore - ne pas générer les clés
	UID     *datastore.Key `json:"uid" datastore:"-"`
	Content string         `json:"content"`
	Author  string         `json:"author"`
}

type Messages struct {
	Messages []Message `json:"messages"`
}

func init() {
	api, _ := endpoints.RegisterService(MessageAPI{}, "messages", "v1", "messages api", true)
	info := api.MethodByName("Add").Info()
	info.Name, info.HTTPMethod, info.Path = "addMessage", "POST", "messages"

	info = api.MethodByName("List").Info()
	info.Name, info.HTTPMethod, info.Path = "getMessages", "GET", "messages"


	endpoints.HandleHTTP()
	// http.HandleFunc("/messages", listMessages)
}

/*
func listMessages(w http.ResponseWriter, r *http.Request) {

	// Récupérer du contexte
	c := appengine.NewContext(r)

	// Création de messages : la liste de message
	messages := []Message{}

	// Récupérer tout les Messages dans le datastore et les mettre dans messages
	_, err := datastore.NewQuery("Message").GetAll(c, &messages)
	// Afficher l'erreur SI il existe
	if err != nil {
		c.Errorf("fetching %v", err)
		return
	}

	// Encoding en json
	enc := json.NewEncoder(w)
	// Récupérer l'erreur
	err = enc.Encode(messages)
	// Afficher l'erreur SI il existe
	if err != nil {
		log.Printf("encoding: %v", err)
		c.Errorf("encoding: %v", err)
	}
}
*/

func (MessageAPI) List(c endpoints.Context) (*Messages, error) {

	messages := []Message{}

	keys, err := datastore.NewQuery("Message").GetAll(c, &messages)
	if err != nil {
		c.Errorf("fetching %v", err)
		return nil, err
	}

	for i, k := range keys {
		messages[i].UID = k
	}

	return &Messages{messages}, nil
}

// Type Message sans les paramètres qui ne nous intéresse pas ID (auto) et Auteur
type AddRequest struct {
	Content string
}

func (MessageAPI) Add(c endpoints.Context, r *AddRequest) (*Message, error) {
	// Récupérer le message
	m := Message{
		Content: r.Content,
	}
	// Récupérer la dernière clé pour Message
	k := datastore.NewIncompleteKey(c, "Message", nil)
	// Ajouter m dans le datastore(c)
	k, err := datastore.Put(c, k, &m)
	// Vérifier que l'erreur n'existe pas
	if err != nil {
		return nil, fmt.Errorf("Put message: %v", err)
	}
	// Récupérer le message ajouter et une erreur nil
	m.UID = k
	return &m, nil
}

/*
type Update struct {
	UID     *datastore.Key
	Content string
}
func (MessageAPI) Update(c endpoints.Context, r *Update) error {
	//
	datastore.RunInTransaction(c, func c appengine.Context() error {

	var message Message

	// Récupérer le message avec son ID
	datastore.GET(c, r.UID, &message)
	// Vérification
	if err != nil {
		return fmt.Errorf("get post: %v", err)
	}
// Mettre
	message.Content = r.Content
	_, err := datastore.Put(c, r.UID, &message)
	})

}
*/
