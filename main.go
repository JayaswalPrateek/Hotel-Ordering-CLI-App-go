package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/inancgumus/screen"
	"github.com/olekukonko/tablewriter"
	"github.com/pterm/pterm"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/exp/maps"
)

type Menu struct {
	Type string // veg or non-veg
	Name string
	Cost int
}
type Bill struct {
	Type string
	Name string
	Cost int
	Qty  int
}

func checkErr(err error, helpInfo string) {
	if err != nil {
		clearScr()
		fmt.Println(helpInfo)
		panic(err)
	}
}
func quit() {
	clearScr()
	fmt.Println("Thanks for visiting, come back later")
	os.Exit(0)
}
func clearScr() {
	screen.Clear()
	screen.MoveTopLeft()
}

func main() {
	clearScr()
	greet()
	menu := decideMenu()
	clearScr()
	printMenu(menu)
	if wannaBuy() {
		bill := takeOrder(menu)
		clearScr()
		fmt.Println("Bill:")
		price := printBill(bill)
		var isItSame bool
		for {
			isItSame, bill = wannaDel(bill)
			if isItSame {
				if askToPay(price) {
					showETA(bill)
					fmt.Println("Order:")
					printBill(bill)
					cookFood(bill)
					thankYou()
				}
				break
			} else {
				clearScr()
				fmt.Println("Bill:")
				price = printBill(bill)
			}
		}
	}
}

func greet() {
	welcome, err := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("W", pterm.NewStyle(pterm.FgLightRed)),
		pterm.NewLettersFromStringWithStyle("E", pterm.NewStyle(pterm.FgYellow)),
		pterm.NewLettersFromStringWithStyle("L", pterm.NewStyle(pterm.FgLightRed)),
		pterm.NewLettersFromStringWithStyle("C", pterm.NewStyle(pterm.FgYellow)),
		pterm.NewLettersFromStringWithStyle("O", pterm.NewStyle(pterm.FgLightRed)),
		pterm.NewLettersFromStringWithStyle("M", pterm.NewStyle(pterm.FgYellow)),
		pterm.NewLettersFromStringWithStyle("E", pterm.NewStyle(pterm.FgLightRed)),
		pterm.NewLettersFromStringWithStyle(" !", pterm.NewStyle(pterm.FgMagenta))).Srender()
	checkErr(err, "pterm caused an internal error")
	pterm.DefaultCenter.Println(welcome)
	pterm.DefaultCenter.WithCenterEachLineSeparately().Printf("🎊️ %s Exclusive Offers 🎊️\n", time.Now().Format("Monday"))
	pterm.DefaultCenter.WithCenterEachLineSeparately().Printf("Enjoy discounts on all dishes every %s !", time.Now().Format("Monday"))
}

func userInputHandler01() int {
	var choice string
	pterm.DefaultBasicText.Printf(pterm.LightBlue("\t~> Your Choice: "))
	fmt.Scanln(&choice)
	choice = strings.TrimSpace(choice)
	choice = strings.ToLower(choice)
	if choice == "" {
		pterm.Error.Printf("Empty input!\n\n")
		return 0
	} else if choice == "exit" {
		quit()
	} else if choice == "yes" {
		return 1
	} else if choice == "no" {
		return 2
	} else {
		pterm.Error.Printf("Invalid input!\n\n")
		return 0
	}
	return 0
}

func decideMenu() map[int]Menu {
	var menu map[int]Menu
	fmt.Println("\n\nAre you Vegan? Answer with yes or no, Enter 'exit' to quit program")
	choice := userInputHandler01()
	decider := true
	for decider {
		if choice != 0 {
			decider = false
			switch choice {
			case 1:
				menu = makeMenuVeg()
			case 2:
				menu = makeMenuNonVeg()
			}
		} else {
			choice = userInputHandler01()
		}
	}
	return menu
}
func makeMenuVeg() map[int]Menu {
	var vegMenu = map[int]Menu{
		1: {"veg", "Burger", 250},
		2: {"veg", "Fries", 100},
		3: {"veg", "Pizza", 450},
		4: {"veg", "Sandwich", 150},
		5: {"veg", "Soup", 300},
		6: {"veg", "Noodles", 350},
		7: {"veg", "Curry", 300},
		8: {"veg", "Rice", 250},
		// 9: {"veg", "Dish Name", Cost},
	}
	return vegMenu
}
func makeMenuNonVeg() map[int]Menu {
	vegMenu := makeMenuVeg()
	vegLen := len(vegMenu)
	var nonVegMenu = map[int]Menu{
		vegLen + 1: {"nonveg", "Chicken Burger", 300},
		vegLen + 2: {"nonveg", "Chicken Wings", 300},
		vegLen + 3: {"nonveg", "Chicken Pizza", 500},
		vegLen + 4: {"nonveg", "Chicken Sandwich", 200},
		// vegLen + 6: {"nonveg", "Dish Name", Cost},
	}
	maps.Copy(nonVegMenu, vegMenu)
	return nonVegMenu
}

func printMenu(menu map[int]Menu) {
	fmt.Println("Menu:")
	column := make([]string, 4)
	allRows := make([][]string, len(menu))
	for i := 1; i <= len(menu); i++ {
		column[0] = strconv.Itoa(i)            // Row numbering
		column[1] = isVeg(menu[i].Type)        // veg or non-veg
		column[2] = menu[i].Name               // dish name
		column[3] = strconv.Itoa(menu[i].Cost) // dish cost
		tmp := make([]string, 4)
		copy(tmp, column)
		allRows = append(allRows, tmp)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"", "", "DISH", "COST(₹)"})
	table.SetHeaderColor(tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold})
	table.SetColumnColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetCenterSeparator("║")
	table.SetColumnSeparator("║")
	table.SetRowSeparator("═")
	table.AppendBulk(allRows)
	table.Render()
}

func isVeg(dishType string) string {
	if dishType == "veg" {
		return "🟢️"
	} else if dishType == "nonveg" {
		return "🔴️"
	} else { // for water, drinks
		return ""
	}
}

func wannaBuy() bool {
	fmt.Println("\nDo you want to order anything from the menu? Answer with yes or no, Enter 'exit' to quit program")
	choice := userInputHandler01()
	decider := true
	for decider {
		if choice != 0 {
			decider = false
			switch choice {
			case 1:
				return true
			case 2:
				quit()
			}
		} else {
			choice = userInputHandler01()
		}
	}
	return false
}

func userInputHandler02(lim int) int {
	var choice string
	pterm.DefaultBasicText.Printf(pterm.LightBlue("\t~> Your Choice: "))
	fmt.Scanln(&choice)
	choice = strings.TrimSpace(choice)
	choice = strings.ToLower(choice)
	if choice == "" {
		pterm.Error.Printf("Empty input!\n\n")
		return 0
	} else if choice == "exit" {
		quit()
	} else if choice == "done" {
		return 1
	} else {
		intChoice, err := strconv.Atoi(choice)
		if err != nil {
			pterm.Error.Printf("Invalid input!\n\n")
			return 0
		} else {
			if intChoice <= 0 {
				pterm.Error.Printf("Invalid input!\n\n")
				return 0
			} else if intChoice > 0 && intChoice <= lim {
				return intChoice + 1234 // workaround, if user entered 1 to select 1st dish, it would trigger "done" switch case, we later subtract 1234
			} else {
				pterm.Error.Printf("Input out of range!\n\n")
				return 0
			}
		}
	}
	return 0
}

func takeOrder(menu map[int]Menu) map[int]Bill {
	bill := make(map[int]Bill)
	var finished, decider, skipFlag bool
	var choice int
	pterm.DefaultBasicText.Printf(pterm.LightMagenta("\nInstructions:\n"))
	pterm.DefaultBasicText.Printf(pterm.LightMagenta("➡️ Enter dish number to add the dish to the cart\n"))
	pterm.DefaultBasicText.Printf(pterm.LightMagenta("➡️ Enter 'done' to checkout\n"))
	pterm.DefaultBasicText.Printf(pterm.LightMagenta("➡️ Enter 'exit' to quit program\n\n"))
	for !finished {
		choice = userInputHandler02(len(menu))
		decider = true
		for decider { // handles mechanism to take a single order
			if choice != 0 {
				decider = false
				skipFlag = false
				switch choice {
				case 1:
					if len(bill) != 0 {
						finished = true
						break
					} else {
						pterm.Error.Printf("No Dishes ordered, instead enter 'exit' to quit\n\n")
					}
				default:
					choice = choice - 1234 // workaround nullified by subtracting 1234
					for i := range bill {  // if one or more items already in the list, check if input is unique, if not then update existing qty
						if bill[i].Name == menu[choice].Name { // if ordered item already in bill, update qty
							if entry, ok := bill[i]; ok {
								entry.Qty++
								bill[i] = entry
							}
							pterm.DefaultBasicText.Println(pterm.LightGreen("\tAdded ", menu[choice].Name, " to cart, now ", bill[i].Qty, " in cart"))
							skipFlag = true // makes sure the repeated item is not relisted in the bill after updating the qty
						}
					}
					if !skipFlag {
						bill[len(bill)+1] = Bill{menu[choice].Type, menu[choice].Name, menu[choice].Cost, 1}
						pterm.DefaultBasicText.Println(pterm.LightGreen("\tAdded ", menu[choice].Name, " to cart"))
					}
				}
			} else {
				choice = userInputHandler02(len(menu))
			}
		}
	}
	for i := range bill {
		if bill[i].Qty != 1 {
			for j := i; j <= len(bill); j++ { // updates cost of item depending upon the qty for all dishes
				if entry, ok := bill[j]; ok {
					entry.Cost = entry.Cost * entry.Qty
					bill[j] = entry
				}
			}
			break
		}
	}
	return bill
}

func printBill(bill map[int]Bill) int {
	column := make([]string, 5)
	allRows := make([][]string, len(bill))
	for i := 1; i <= len(bill); i++ {
		column[0] = strconv.Itoa(i)
		column[1] = isVeg(bill[i].Type)
		column[2] = bill[i].Name
		column[3] = strconv.Itoa(bill[i].Qty)
		column[4] = strconv.Itoa(bill[i].Cost)
		tmp := make([]string, 5)
		copy(tmp, column)
		allRows = append(allRows, tmp)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"", "", "DISH", "Qty", "COST(₹)"})
	table.SetHeaderColor(tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold})
	table.AppendBulk(allRows)
	table.SetColumnColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetCenterSeparator("║")
	table.SetColumnSeparator("║")
	table.SetRowSeparator("═")
	table.Render()
	price := showTotal(bill)
	return price
}
func showTotal(bill map[int]Bill) int {
	var total int
	for i := range bill {
		total += bill[i].Cost
	}
	fmt.Printf("\n")
	panel := pterm.DefaultBox.Sprint("Total Cost: ₹ ", strconv.Itoa(total))
	pterm.Println(panel)
	return total
}

func wannaDel(bill map[int]Bill) (bool, map[int]Bill) {
	fmt.Println("\nDo you want to remove any items? Answer with yes or no, Enter 'exit' to quit program")
	choice := userInputHandler01()
	decider := true
	for decider {
		if choice != 0 {
			decider = false
			switch choice {
			case 1:
				bill = letsDel(bill)
				return false, bill
			case 2:
				return true, bill
			}
		} else {
			choice = userInputHandler01()
		}
	}
	return false, bill
}
func letsDel(bill map[int]Bill) map[int]Bill {
	var finished, decider, wasRemoved, skipLoop bool
	var history []int
	var choice, unitCost int
	billLen := len(bill)
	pterm.DefaultBasicText.Printf(pterm.LightMagenta("\nInstructions:\n"))
	pterm.DefaultBasicText.Printf(pterm.LightMagenta("➡️ Enter dish number to remove the dish from the cart\n"))
	pterm.DefaultBasicText.Printf(pterm.LightMagenta("➡️ Enter 'done' to confirm\n"))
	pterm.DefaultBasicText.Printf(pterm.LightMagenta("➡️ Enter 'exit' to quit program\n\n"))
	for !finished {
		choice = userInputHandler02(billLen)
		copyChoice := choice - 1234
		skipLoop = false
		for _, val := range history {
			if copyChoice == val {
				pterm.Error.Printf("dish %d was already removed from bill, Enter 'done' to print the new bill\n\n", copyChoice)
				skipLoop = true
				break
			}
		}
		if skipLoop {
			continue
		}
		decider = true
		for decider {
			if choice != 0 {
				decider = false
				switch choice {
				case 1:
					if len(bill) != 0 {
						history = nil
						finished = true
						break
					} else {
						clearScr()
						fmt.Println("😭️ All Items Removed from cart 😭️")
						fmt.Println("Thanks for visiting, come back later")
						os.Exit(0)
					}
				default:
					choice = choice - 1234
					if bill[choice].Qty == 1 {
						pterm.DefaultBasicText.Println(pterm.LightGreen("deleted ", bill[choice].Name, " from cart"))
						history = append(history, choice)
						delete(bill, choice)
						wasRemoved = true
					} else if bill[choice].Qty > 1 {
						var howMany int
						for {
							fmt.Printf("There are %dx %s in cart, how many to remove: ", bill[choice].Qty, bill[choice].Name)
							fmt.Scanln(&howMany)
							if howMany == 0 {
								fmt.Printf("Number of %s was not modified\n\n", bill[choice].Name)
								break
							} else if howMany < 0 {
								fmt.Println("Invalid input!")
							} else if howMany < 0 || howMany > bill[choice].Qty {
								pterm.Error.Printf("Input out of range\n\n")
								continue
							} else if howMany == bill[choice].Qty {
								pterm.DefaultBasicText.Println(pterm.LightGreen("deleted ", bill[choice].Name, " from cart"))
								history = append(history, choice)
								delete(bill, choice)
								wasRemoved = true
								break
							} else {
								if entry, ok := bill[choice]; ok {
									entry.Qty -= howMany
									// need to update cost after quantity is changed
									unitCost = bill[choice].Cost / bill[choice].Qty
									entry.Cost -= unitCost * howMany
									bill[choice] = entry
									wasRemoved = false
									pterm.DefaultBasicText.Printf(pterm.LightGreen(bill[choice].Qty, "x ", bill[choice].Name, " now left in cart\n\n"))
									break
								} else {
									pterm.Error.Println("Internal Error")
								}
							}
						}
					}
				}
			} else {
				choice = userInputHandler02(billLen)
			}
		}
	}
	if wasRemoved { // reorder items in the bill if any item was removed
		return resetKeys(bill)
	} else {
		return bill
	}
}
func resetKeys(bill map[int]Bill) map[int]Bill { // makes sure that empty bill does not have empty entry after deletion
	newBill, ctr := make(map[int]Bill), 1
	for key, val := range bill {
		newBill[ctr] = bill[key]
		newBill[ctr] = val
		ctr++
	}
	return newBill
}

func askToPay(price int) bool {
	fmt.Printf("\nProceed to pay ₹ ")
	pterm.NewRGB(133, 187, 101).Printf(strconv.Itoa(price))
	fmt.Println(" ? Answer with yes or no, Enter 'exit' to quit program")
	choice := userInputHandler01()
	decider := true
	for decider {
		if choice != 0 {
			decider = false
			switch choice {
			case 1:
				clearScr()
				pterm.Success.Println("Transaction Successful")
				return true
			case 2:
				quit()
			}
		} else {
			choice = userInputHandler01()
		}
	}
	return false
}

func showETA(bill map[int]Bill) {
	var username string
	var eta int
	fmt.Printf("\nPlease Enter your name: ")
	fmt.Scanln(&username)
	clearScr()
	for i := range bill {
		eta = eta + 4*bill[i].Qty // 4 seconds to make each dish
	}
	if eta >= 60 {
		min := eta / 60
		sec := eta % 60
		fmt.Printf("Currently preparing for: %s, ETA is %d min and %d sec\n\n", username, min, sec)
	} else {
		fmt.Printf("Currently preparing for: %s, ETA is %d sec \n\n", username, eta)
	}
}
func cookFood(bill map[int]Bill) {
	fmt.Printf("\n\n")
	for i := range bill {
		fmt.Printf("Preparing %s now...\n", bill[i].Name)
		if bill[i].Qty == 1 {
			time.Sleep(4 * time.Second)
		}
		showBar(bill[i].Qty)
		fmt.Printf("\n")
		fmt.Printf("       %s is ready !\n\n", bill[i].Name)
	}
}
func showBar(qty int) {
	bar := progressbar.NewOptions(qty,
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionShowCount(),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionShowBytes(false),
	)
	for i := 0; i < qty; i++ {
		err := bar.Add(1)
		checkErr(err, "progress bar caused an internal error")
		if qty == 1 {
			time.Sleep(0 * time.Second)
		} else {
			time.Sleep(4 * time.Second) // 4 seconds to make each dish
		}
	}
}
func thankYou() {
	err := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("THANK YOU !", pterm.NewStyle(pterm.FgMagenta))).Render()
	checkErr(err, "pterm caused an internal error")
}
