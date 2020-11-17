package structs

type Status struct {
	Success bool
}

type Request struct {
	FunctionName string
}

type Location struct {
	ID      int
	Name    string
	MenuID  int
	Address string
	Hours   string
	Phone   int
}

type Menu struct {
	ID         int
	Name       string
	LocationID int
	Items      []FoodItem
}

type FoodItem struct {
	ID             int
	Name           string
	Description    string
	Cost           int
	IsAvailable    bool
	Ingredients    string
	NutritionFacts string
}

type User struct {
	ID               int
	DollarBalance    int
	MealSwipeBalance int
}

type ItemOrder struct {
	ID            int
	ItemID        int
	Customization string
}

type Order struct {
	ID 		   int
	UserID     int
	LocationID int
	Item       ItemOrder
	Status     string
	SubmitTime string
}

type OrderAndItems struct {
	Order Order
	Items []ItemOrder
}

type ListLocationsResponse struct {
	Locations []Location
}

type GetMenuRequest struct {
	MenuID int
}

type GetMenuResponse struct {
	RequestedMenu Menu
}

type ViewItemRequest struct {
	ItemID int
}

type ViewItemResponse struct {
	Item FoodItem
}

type CreateOrderRequest struct {
	OrderRequest Order
	PayWithMealSwipe bool
}

type CreateOrderResponse struct {
	ID      int
	Success bool
	Status  string
}

type CheckOrderStatusRequest struct {
	OrderID int
}

type CheckOrderStatusResponse struct {
	Status string
}

type GetOrdersRequest struct {
	LocationID int
}

type GetOrdersResponse struct {
	Orders []ItemOrder
}

type SelectOrderRequest struct {
	OrderID int
}

type SelectOrderResponse struct {
	Status string
}

type CompelteOrderRequest struct {
	OrderID int
}

type CompelteOrderResponse struct {
	Status string
}

type UpdateItemRequest struct {
	MenuID  int
	ItemID  int
	NewItem FoodItem
}

type UpdateItemResponse struct {
	Status string
}

type CreateItemRequest struct {
	MenuID  int
	NewItem FoodItem
}

type CreateItemResponse struct {
	Status string
}

type DeleteItemRequest struct {
	MenuID int
	ItemID int
}

type DeleteItemResponse struct {
	Status string
}

type UpdateOrderRequest struct {
	ID int
}

type AddItemToOrderRequest struct {
	OrderID int
	Item ItemOrder
	PayWithMealSwipe bool
}

type AddItemToOrderResponse struct {

}

type SendMealSwipesRequest struct {
	FromID int
	ToID int
	NumSwipes int
}

type SendMealSwipesResponse struct {
	Success bool
	Balance int 
}

type GetPaymentBalancesRequest struct {
	UserID int
}

type GetPaymentBalancesResponse struct {
	MealSwipeBalance int
	CentsBalance int
}