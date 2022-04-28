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
	fmt.Println("Thanks for visiting us, do come back later 😚️")
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
					break
				}
			} else {
				clearScr()
				fmt.Println("Bill:")
				price = printBill(bill)
			}
		}
	}
}

func greet() {
	banner, err := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("W", pterm.NewStyle(pterm.FgLightRed)),
		pterm.NewLettersFromStringWithStyle("E", pterm.NewStyle(pterm.FgYellow)),
		pterm.NewLettersFromStringWithStyle("L", pterm.NewStyle(pterm.FgLightRed)),
		pterm.NewLettersFromStringWithStyle("C", pterm.NewStyle(pterm.FgYellow)),
		pterm.NewLettersFromStringWithStyle("O", pterm.NewStyle(pterm.FgLightRed)),
		pterm.NewLettersFromStringWithStyle("M", pterm.NewStyle(pterm.FgYellow)),
		pterm.NewLettersFromStringWithStyle("E", pterm.NewStyle(pterm.FgLightRed)),
		pterm.NewLettersFromStringWithStyle(" !", pterm.NewStyle(pterm.FgMagenta))).Srender()
	checkErr(err, "pterm caused an internal error")
	pterm.DefaultCenter.Println(banner)
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
		fmt.Printf("\t")
		pterm.Error.Printf("Empty input!\n\n")
		return 0
	} else if choice == "exit" {
		quit()
	} else if choice == "yes" {
		return 1
	} else if choice == "no" {
		return 2
	} else {
		fmt.Printf("\t")
		pterm.Error.Printf("Invalid input!\n\n")
		return 0
	}
	return 0
}

func decideMenu() map[int]Menu {
	var menu map[int]Menu
	var stopper bool
	fmt.Println("\n\nAre you Vegan? Answer with yes or no, Enter 'exit' to quit")
	choice := userInputHandler01()
	for !stopper {
		if choice != 0 {
			stopper = true
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
		// 9: {"veg", "Name", Cost},
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
		// vegLen + 5: {"nonveg", "Dish Name", Cost},
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
		return "  "
	}
}

func wannaBuy() bool {
	fmt.Println("\nDo you want to order anything from our menu? Answer with yes or no, Enter 'exit' to quit")
	choice := userInputHandler01()
	var stopper bool
	for !stopper {
		if choice != 0 {
			stopper = true
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
		fmt.Printf("\t")
		pterm.Error.Printf("Empty input!\n\n")
		return 0
	} else if choice == "exit" {
		quit()
	} else if choice == "done" {
		return 1
	} else {
		intChoice, err := strconv.Atoi(choice)
		if err != nil {
			fmt.Printf("\t")
			pterm.Error.Printf("Invalid input!\n\n")
			return 0
		} else {
			if intChoice <= 0 {
				fmt.Printf("\t")
				pterm.Error.Printf("Invalid input!\n\n")
				return 0
			} else if intChoice > 0 && intChoice <= lim {
				return intChoice + 1234 // workaround, if user entered 1 to select 1st dish, it would trigger "done" in the switch case, we later subtract 1234
			} else {
				fmt.Printf("\t")
				pterm.Error.Printf("Input out of range!\n\n")
				return 0
			}
		}
	}
	return 0
}

func userInputHandler03() int {
	var choice string
	fmt.Scanln(&choice)
	choice = strings.TrimSpace(choice)
	choice = strings.ToLower(choice)
	if choice == "" {
		fmt.Printf("\t\t")
		pterm.Error.Printf("Empty input!\n\n")
		return 0
	} else {
		intChoice, err := strconv.Atoi(choice)
		if err != nil {
			fmt.Printf("\t\t")
			pterm.Error.Printf("Invalid input!\n\n")
			return 0
		} else {
			if intChoice <= 0 {
				fmt.Printf("\t\t")
				pterm.Error.Printf("Invalid input!\n\n")
				return 0
			} else {
				return intChoice
			}
		}
	}
}
func takeOrder(menu map[int]Menu) map[int]Bill {
	bill := make(map[int]Bill)
	var finished, stopper, skipFlag bool
	var choice int
	pterm.DefaultBasicText.Printf(pterm.LightMagenta("\nInstructions:\n"))
	pterm.DefaultBasicText.Printf(pterm.LightMagenta("➡️ Enter dish number to add the dish to the cart\n"))
	pterm.DefaultBasicText.Printf(pterm.LightMagenta("➡️ Enter 'done' to checkout\n"))
	pterm.DefaultBasicText.Printf(pterm.LightMagenta("➡️ Enter 'exit' to quit\n\n"))
	for !finished {
		choice = userInputHandler02(len(menu))
		stopper = false // reset value
		for !stopper {
			if choice != 0 {
				stopper = true
				skipFlag = false // reset value
				switch choice {
				case 1:
					if len(bill) != 0 {
						finished = true
						break
					} else {
						pterm.Error.Printf("No Dishes ordered, instead enter 'exit' to quit\n\n")
					}
				default:
					var howMany int
					choice = choice - 1234 // workaround nullified here
					for {
						fmt.Printf("\t\tHow many %s do you want to order: ", menu[choice].Name)
						howMany = userInputHandler03()
						if howMany != 0 {
							break
						}
					}
					for i := range bill {
						if bill[i].Name == menu[choice].Name { // check if input dish is already in bill, if yes then update its existing qty and cost
							if entry, ok := bill[i]; ok {
								entry.Qty += howMany                       // update qty
								entry.Cost = entry.Qty * menu[choice].Cost // update cost
								bill[i] = entry
							}
							pterm.DefaultBasicText.Println(pterm.LightGreen("Added ", howMany, " ", menu[choice].Name, " to cart, now ", bill[i].Qty, " in cart"))
							skipFlag = true // makes sure the repeated item is not relisted in the bill after updating the qty in the bill
						}
					}
					if !skipFlag {
						bill[len(bill)+1] = Bill{menu[choice].Type, menu[choice].Name, menu[choice].Cost * howMany, howMany}
						pterm.DefaultBasicText.Println(pterm.LightGreen("Added ", howMany, " ", menu[choice].Name, " to cart"))
					}
				}
			} else {
				choice = userInputHandler02(len(menu))
			}
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
	return showTotal(bill)
}
func showTotal(bill map[int]Bill) int {
	var total int
	for i := range bill {
		total += bill[i].Cost
	}
	fmt.Printf("\n")
	pterm.Println(pterm.DefaultBox.Sprint("Total Cost: ₹ ", strconv.Itoa(total)))
	return total
}

func userInputHandler04(lim int) int {
	var choice string
	fmt.Scanln(&choice)
	choice = strings.TrimSpace(choice)
	choice = strings.ToLower(choice)
	if choice == "" {
		fmt.Printf("\t\t")
		pterm.Error.Printf("Empty input!\n\n")
		return 0
	} else {
		intChoice, err := strconv.Atoi(choice)
		if err != nil {
			fmt.Printf("\t\t")
			pterm.Error.Printf("Invalid input!\n\n")
			return 0
		} else {
			if intChoice < 0 {
				fmt.Printf("\t\t")
				pterm.Error.Printf("Invalid input!\n\n")
				return 0
			} else if intChoice >= 0 && intChoice <= lim {
				return intChoice + 1234 // workaround, if user entered 1 to select 1st dish, it would trigger "done" in the switch case, we later subtract 1234
			} else {
				fmt.Printf("\t\t")
				pterm.Error.Printf("Input out of range!\n\n")
				return 0
			}
		}
	}
}
func wannaDel(bill map[int]Bill) (bool, map[int]Bill) {
	var stopper bool
	fmt.Println("\nDo you want to remove any items? Answer with yes or no, Enter 'exit' to quit")
	choice := userInputHandler01()
	for !stopper {
		if choice != 0 {
			stopper = true
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
	var finished, stopper, wasRemoved, skipLoop bool
	var history []int
	var choice, unitCost, howMany int
	billLen := len(bill)
	pterm.DefaultBasicText.Printf(pterm.LightMagenta("\nInstructions:\n"))
	pterm.DefaultBasicText.Printf(pterm.LightMagenta("➡️ Enter dish number to remove the dish from the cart\n"))
	pterm.DefaultBasicText.Printf(pterm.LightMagenta("➡️ Enter 'done' to proceed\n"))
	pterm.DefaultBasicText.Printf(pterm.LightMagenta("➡️ Enter 'exit' to quit\n\n"))
	for !finished {
		choice = userInputHandler02(billLen)
		copyChoice := choice - 1234 // workaround nullified here
		skipLoop = false            // reset value
		for _, removedItem := range history {
			if copyChoice == removedItem {
				pterm.Error.Printf("dish %d was already removed from bill, Enter 'done' to print the new bill\n\n", copyChoice)
				skipLoop = true
				break
			}
		}
		if skipLoop {
			continue
		}
		stopper = false // reset value
		for !stopper {
			if choice != 0 {
				stopper = true
				switch choice {
				case 1:
					if len(bill) != 0 {
						history = nil
						finished = true
						break
					} else {
						clearScr()
						fmt.Println("😭️ All Items Removed from cart 😭️")
						fmt.Println("Thanks for visiting us, do come back later 😚️")
						os.Exit(0)
					}
				default:
					choice = choice - 1234
					if bill[choice].Qty == 1 {
						pterm.DefaultBasicText.Println(pterm.LightGreen("deleted ", bill[choice].Name, " from cart"))
						history = append(history, choice)
						delete(bill, choice)
						wasRemoved = true
					} else {
						var stopper2 bool
						for !stopper2 {
							fmt.Printf("\t\tThere are %dx %s in cart, enter how many to remove: ", bill[choice].Qty, bill[choice].Name)
							howMany = userInputHandler04(bill[choice].Qty)
							if howMany != 0 {
								stopper2 = true
								howMany -= 1234
								if howMany == 0 {
									fmt.Printf("\t\t")
									pterm.Info.Println("Number of", bill[choice].Name, "was not modified")
									fmt.Printf("\n")
									break
								} else if howMany == bill[choice].Qty {
									pterm.DefaultBasicText.Println(pterm.LightGreen("deleted ", bill[choice].Name, " from cart"))
									history = append(history, choice)
									delete(bill, choice)
									wasRemoved = true
									break
								} else {
									if entry, ok := bill[choice]; ok {
										unitCost = entry.Cost / entry.Qty
										entry.Cost -= unitCost * howMany
										entry.Qty -= howMany
										bill[choice] = entry
										pterm.DefaultBasicText.Printf(pterm.LightGreen(bill[choice].Qty, "x ", bill[choice].Name, " now left in cart\n\n"))
										break
									}
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
	if wasRemoved {
		return resetKeys(bill)
	} else {
		return bill
	}
}
func resetKeys(bill map[int]Bill) map[int]Bill {
	newBill, ctr := make(map[int]Bill), 1
	for key, val := range bill {
		newBill[ctr] = bill[key]
		newBill[ctr] = val
		ctr++
	}
	return newBill
}

func askToPay(price int) bool {
	var stopper bool
	fmt.Printf("\nProceed to pay ₹ ")
	pterm.NewRGB(133, 187, 101).Printf(strconv.Itoa(price))
	fmt.Println(" ? Answer with yes or no, Enter 'exit' to quit program")
	choice := userInputHandler01()
	for !stopper {
		if choice != 0 {
			stopper = true
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
	var err error
	fmt.Printf("\n\n")
	for i := range bill {
		fmt.Printf("Preparing %s now...\n", bill[i].Name)
		bar := progressbar.NewOptions(bill[i].Qty,
			progressbar.OptionSetRenderBlankState(true),
			progressbar.OptionFullWidth(),
			progressbar.OptionShowCount(),
			progressbar.OptionSetPredictTime(false),
			progressbar.OptionShowBytes(false),
		)
		for j := 1; j <= bill[i].Qty; j++ {
			time.Sleep(4 * time.Second)
			err = bar.Add(1)
			checkErr(err, "progress bar caused an internal error")
		}
		fmt.Printf("\n")
		fmt.Printf("       %s is ready !\n\n", bill[i].Name)
	}
}
func thankYou() {
	err := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("THANK YOU !", pterm.NewStyle(pterm.FgMagenta))).Render()
	checkErr(err, "pterm caused an internal error")
}
