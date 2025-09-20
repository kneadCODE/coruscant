package entity

import (
	"fmt"

	"github.com/kneadCODE/coruscant/shared/golib/id"
)

// ItemID represents a unique identifier for a budget item using UUIDv7
type ItemID struct {
	id.EntityID
}

// NewItemID creates a new ItemID using UUIDv7
func NewItemID() (ItemID, error) {
	base, err := id.NewEntityID()
	if err != nil {
		return ItemID{}, fmt.Errorf("failed to create item ID: %w", err)
	}
	return ItemID{EntityID: base}, nil
}

// NewItemIDFromString creates an ItemID from an existing string
func NewItemIDFromString(idStr string) (ItemID, error) {
	base, err := id.NewEntityIDFromString(idStr)
	if err != nil {
		return ItemID{}, fmt.Errorf("failed to create item ID: %w", err)
	}
	return ItemID{EntityID: base}, nil
}

// Equals checks if two ItemIDs are equal
func (i ItemID) Equals(other ItemID) bool {
	return i.EntityID.Equals(other.EntityID)
}

// ItemType represents the type of budget item
type ItemType string

// Item type constants define the type of budget items
const (
	ItemTypeIncome   ItemType = "INCOME"
	ItemTypeExpense  ItemType = "EXPENSE"
	ItemTypeTransfer ItemType = "TRANSFER"
)

// NewItemType creates a new ItemType from string
func NewItemType(itemType string) (ItemType, error) {
	switch ItemType(itemType) {
	case ItemTypeIncome, ItemTypeExpense, ItemTypeTransfer:
		return ItemType(itemType), nil
	default:
		return ItemTypeExpense, fmt.Errorf("invalid item type: %s", itemType)
	}
}

// String returns the string representation of ItemType
func (i ItemType) String() string {
	return string(i)
}

// IsIncome checks if the item type is income
func (i ItemType) IsIncome() bool {
	return i == ItemTypeIncome
}

// IsExpense checks if the item type is expense
func (i ItemType) IsExpense() bool {
	return i == ItemTypeExpense
}

// IsTransfer checks if the item type is transfer
func (i ItemType) IsTransfer() bool {
	return i == ItemTypeTransfer
}

// ItemCategory represents the category of an item
type ItemCategory string

// Income categories
const (
	ItemCategorySalary            ItemCategory = "SALARY"
	ItemCategoryFreelance         ItemCategory = "FREELANCE"
	ItemCategoryInvestmentReturns ItemCategory = "INVESTMENT_RETURNS"
	ItemCategoryRentalIncome      ItemCategory = "RENTAL_INCOME"
	ItemCategorySideBusiness      ItemCategory = "SIDE_BUSINESS"
	ItemCategoryBonus             ItemCategory = "BONUS"
	ItemCategoryGifts             ItemCategory = "GIFTS"
	ItemCategoryOtherIncome       ItemCategory = "OTHER_INCOME"
)

// Expense categories
const (
	ItemCategoryHousing        ItemCategory = "HOUSING"
	ItemCategoryFood           ItemCategory = "FOOD"
	ItemCategoryTransportation ItemCategory = "TRANSPORTATION"
	ItemCategoryUtilities      ItemCategory = "UTILITIES"
	ItemCategoryInsurance      ItemCategory = "INSURANCE"
	ItemCategoryHealthcare     ItemCategory = "HEALTHCARE"
	ItemCategoryEntertainment  ItemCategory = "ENTERTAINMENT"
	ItemCategoryEducation      ItemCategory = "EDUCATION"
	ItemCategoryShopping       ItemCategory = "SHOPPING"
	ItemCategoryPersonalCare   ItemCategory = "PERSONAL_CARE"
	ItemCategorySubscriptions  ItemCategory = "SUBSCRIPTIONS"
	ItemCategoryTaxes          ItemCategory = "TAXES"
	ItemCategoryOtherExpense   ItemCategory = "OTHER_EXPENSE"
)

// Transfer categories
const (
	ItemCategorySavingsTransfer        ItemCategory = "SAVINGS_TRANSFER"
	ItemCategoryDebtPayment            ItemCategory = "DEBT_PAYMENT"
	ItemCategoryInvestmentContribution ItemCategory = "INVESTMENT_CONTRIBUTION"
	ItemCategoryEmergencyFund          ItemCategory = "EMERGENCY_FUND"
	ItemCategoryRetirementContribution ItemCategory = "RETIREMENT_CONTRIBUTION"
	ItemCategoryOtherTransfer          ItemCategory = "OTHER_TRANSFER"
)

// itemCategoryValues maps string values to ItemCategory constants for validation
var itemCategoryValues = map[string]ItemCategory{
	"SALARY":                  ItemCategorySalary,
	"FREELANCE":               ItemCategoryFreelance,
	"INVESTMENT_RETURNS":      ItemCategoryInvestmentReturns,
	"RENTAL_INCOME":           ItemCategoryRentalIncome,
	"SIDE_BUSINESS":           ItemCategorySideBusiness,
	"BONUS":                   ItemCategoryBonus,
	"GIFTS":                   ItemCategoryGifts,
	"OTHER_INCOME":            ItemCategoryOtherIncome,
	"HOUSING":                 ItemCategoryHousing,
	"FOOD":                    ItemCategoryFood,
	"TRANSPORTATION":          ItemCategoryTransportation,
	"UTILITIES":               ItemCategoryUtilities,
	"INSURANCE":               ItemCategoryInsurance,
	"HEALTHCARE":              ItemCategoryHealthcare,
	"ENTERTAINMENT":           ItemCategoryEntertainment,
	"EDUCATION":               ItemCategoryEducation,
	"SHOPPING":                ItemCategoryShopping,
	"PERSONAL_CARE":           ItemCategoryPersonalCare,
	"SUBSCRIPTIONS":           ItemCategorySubscriptions,
	"TAXES":                   ItemCategoryTaxes,
	"OTHER_EXPENSE":           ItemCategoryOtherExpense,
	"SAVINGS_TRANSFER":        ItemCategorySavingsTransfer,
	"DEBT_PAYMENT":            ItemCategoryDebtPayment,
	"INVESTMENT_CONTRIBUTION": ItemCategoryInvestmentContribution,
	"EMERGENCY_FUND":          ItemCategoryEmergencyFund,
	"RETIREMENT_CONTRIBUTION": ItemCategoryRetirementContribution,
	"OTHER_TRANSFER":          ItemCategoryOtherTransfer,
}

// NewItemCategory creates a new ItemCategory from string
func NewItemCategory(category string) (ItemCategory, error) {
	if c, exists := itemCategoryValues[category]; exists {
		return c, nil
	}
	return ItemCategoryOtherExpense, fmt.Errorf("invalid item category: %s", category)
}

// String returns the string representation of ItemCategory
func (i ItemCategory) String() string {
	return string(i)
}

// GetItemType returns the ItemType for this category
func (i ItemCategory) GetItemType() ItemType {
	switch i {
	case ItemCategorySalary, ItemCategoryFreelance, ItemCategoryInvestmentReturns,
		ItemCategoryRentalIncome, ItemCategorySideBusiness, ItemCategoryBonus,
		ItemCategoryGifts, ItemCategoryOtherIncome:
		return ItemTypeIncome
	case ItemCategorySavingsTransfer, ItemCategoryDebtPayment, ItemCategoryInvestmentContribution,
		ItemCategoryEmergencyFund, ItemCategoryRetirementContribution, ItemCategoryOtherTransfer:
		return ItemTypeTransfer
	default:
		return ItemTypeExpense
	}
}

// GetIncomeCategories returns all income categories
func GetIncomeCategories() []ItemCategory {
	return []ItemCategory{
		ItemCategorySalary,
		ItemCategoryFreelance,
		ItemCategoryInvestmentReturns,
		ItemCategoryRentalIncome,
		ItemCategorySideBusiness,
		ItemCategoryBonus,
		ItemCategoryGifts,
		ItemCategoryOtherIncome,
	}
}

// GetExpenseCategories returns all expense categories
func GetExpenseCategories() []ItemCategory {
	return []ItemCategory{
		ItemCategoryHousing,
		ItemCategoryFood,
		ItemCategoryTransportation,
		ItemCategoryUtilities,
		ItemCategoryInsurance,
		ItemCategoryHealthcare,
		ItemCategoryEntertainment,
		ItemCategoryEducation,
		ItemCategoryShopping,
		ItemCategoryPersonalCare,
		ItemCategorySubscriptions,
		ItemCategoryTaxes,
		ItemCategoryOtherExpense,
	}
}

// GetTransferCategories returns all transfer categories
func GetTransferCategories() []ItemCategory {
	return []ItemCategory{
		ItemCategorySavingsTransfer,
		ItemCategoryDebtPayment,
		ItemCategoryInvestmentContribution,
		ItemCategoryEmergencyFund,
		ItemCategoryRetirementContribution,
		ItemCategoryOtherTransfer,
	}
}
