package structs

type Status struct {
	Success bool
}

type Request struct {
	FunctionName string
	Data         interface{}
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
	NutritionFacts string
}

type User struct {
	ID               int
	DollarBalance    int
	MealSwipeBalance int
}

type OrderItem struct {
	ID            int
	FoodID        int
	Customization string
	PayWithSwipe  bool
}

type OrderItemWithFood struct {
	Item OrderItem
	Food FoodItem
}

type Order struct {
	ID               int
	UserID           int
	LocationID       int
	Item             OrderItem
	Status           string
	SubmitTime       string
	LastStatusChange string
	SwipeCost        int
	CentCost         int
}

type OrderAndItems struct {
	Order Order
	Items []OrderItem
}

type OrderAndItemsWithFood struct {
	Order Order
	Items []OrderItemWithFood
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
}

type CreateOrderResponse struct {
	ID      int
	Success bool
	Status  string
}

type GetOrderHistoryRequest struct {
	UserID int
}

type GetOrdersRequest struct {
	LocationID int
}

type GetOrdersResponse struct {
	Orders []OrderItem
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
	Item    OrderItem
}

type AddItemToOrderResponse struct {
}

type DeleteItemFromOrderRequest struct {
	ItemID  int
	OrderID int
}

type SendMealSwipesRequest struct {
	FromID    int
	ToID      int
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
	CentsBalance     int
}

type GetCartRequest struct {
	UserID int
}

type LoginRequest struct {
	UserNetID string
	Password  string
}

type LoginResponse struct {
	Status    bool
	IsStudent bool
	UserID    int
}
