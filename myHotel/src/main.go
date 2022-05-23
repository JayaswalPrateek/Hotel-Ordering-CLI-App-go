package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/olekukonko/tablewriter"
	"github.com/pterm/pterm"
	"github.com/schollz/progressbar/v3"
)

func main() {
	clearScr()
	greet()
	menu := decideMenu()
	clearScr()
	printMenu(menu)
	if wannaBuy() {
		clearScr()
		printMenu(menu)
		bill := takeOrder(menu)
		clearScr()
		price := printBill(bill)
		var isItSame bool
		for {
			isItSame, bill = wannaDel(bill)
			if isItSame {
				clearScr()
				printBill(bill)
				if askToPay(price) {
					showETA(bill)
					printOrder(bill)
					cookFood(bill)
					thankYou()
					showCursor()
					break
				}
			} else {
				clearScr()
				price = printBill(bill)
			}
		}
	}
}

type Menu struct {
	Type string
	Name string
	Cost int
}
type Bill struct {
	Type string
	Name string
	Cost int
	Qty  int
}

func hideCursor() {
	fmt.Print("\033[?25l")
}
func showCursor() {
	fmt.Print("\033[?25h")
}
func clearScr() {
	fmt.Print("\x1bc")
}
func quit() {
	clearScr()
	pterm.DefaultCenter.Print(pterm.Italic.Sprint("Thanks for visiting us, do come back later 😚️"))
	showCursor()
	os.Exit(0)
}
func checkErr(err error, helpInfo string) {
	if err != nil {
		clearScr()
		fmt.Println(helpInfo)
		showCursor()
		panic(err)
	}
}
func isVeg(checkDish string) string {
	switch checkDish {
	case "veg":
		return "🟢️"
	case "nonveg":
		return "🔴️"
	default:
		return "  "
	}
}
func fetchVar01() string {
	hideCursor()
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	checkErr(err, "bufio usage caused an error")
	return strings.TrimSpace(choice)
}

func greet() {
	banner, err := (pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("W", pterm.NewStyle(pterm.FgLightRed)),
		pterm.NewLettersFromStringWithStyle("E", pterm.NewStyle(pterm.FgYellow)),
		pterm.NewLettersFromStringWithStyle("L", pterm.NewStyle(pterm.FgLightRed)),
		pterm.NewLettersFromStringWithStyle("C", pterm.NewStyle(pterm.FgYellow)),
		pterm.NewLettersFromStringWithStyle("O", pterm.NewStyle(pterm.FgLightRed)),
		pterm.NewLettersFromStringWithStyle("M", pterm.NewStyle(pterm.FgYellow)),
		pterm.NewLettersFromStringWithStyle("E", pterm.NewStyle(pterm.FgLightRed)),
		pterm.NewLettersFromStringWithStyle(" !", pterm.NewStyle(pterm.FgMagenta))).Srender())
	checkErr(err, "pterm usage caused an error")
	pterm.DefaultCenter.Println(banner)
	pterm.DefaultCenter.Printfln("🎊️ %s Exclusive Offers 🎊️", time.Now().Format("Monday"))
	pterm.DefaultCenter.Printfln("Enjoy discounts on all dishes every %s !", time.Now().Format("Monday"))
}

func userInputHandler01() int {
	pterm.NewStyle(pterm.Bold, pterm.FgMagenta).Print("~> ")
	choice := strings.ToLower(fetchVar01())
	switch choice {
	case "yes":
		return 1
	case "no":
		return 2
	case "exit":
		quit()
	case "":
		fmt.Print(pterm.Error.Sprint("Empty input!"), "\n\n")
		return 0
	default:
		fmt.Print(pterm.Error.Sprint("Invalid input!"), "\n\n")
		return 0
	}
	return 0
}

func decideMenu() map[int]Menu {
	var stopper bool
	pterm.DefaultCenter.Print(pterm.DefaultBox.Sprint(pterm.NewStyle(pterm.Bold, pterm.FgCyan).Sprint("Are you Vegan? Answer with 'yes' or 'no', Enter 'exit' to quit")))
	for !stopper {
		choice := userInputHandler01()
		if choice == 0 {
			continue
		}
		switch choice {
		case 1:
			return setMenuVeg()
		case 2:
			return setMenuNonVeg()
		}
	}
	return nil
}
func setMenuVeg() map[int]Menu {
	var vegMenu = map[int]Menu{
		1: {"veg", "Burger", 250},
		2: {"veg", "Fries", 100},
		3: {"veg", "Pizza", 350},
		4: {"veg", "Sandwich", 150},
		5: {"veg", "Soup", 200},
		6: {"veg", "Noodles", 350},
		7: {"veg", "Curry", 300},
		8: {"veg", "Rice", 250},
		// 9: {"veg", "Name", Cost},
	}
	return vegMenu
}
func setMenuNonVeg() map[int]Menu {
	vegMenu := setMenuVeg()
	vegLen := len(vegMenu)
	var nonVegMenu = map[int]Menu{
		vegLen + 1: {"nonveg", "Chk-Burger", 300},
		vegLen + 2: {"nonveg", "Chk-Wings", 200},
		vegLen + 3: {"nonveg", "Chk-Pizza", 400},
		vegLen + 4: {"nonveg", "Chk-Sandwich", 200},
		vegLen + 5: {"nonveg", "Chk-Curry", 350},
		// vegLen + 6: {"nonveg", "Dish Name", Cost},
	}
	err := mergo.Merge(&nonVegMenu, vegMenu)
	checkErr(err, "mergo usage caused an error")
	return nonVegMenu
}
func printMenu(menu map[int]Menu) {
	allRows := make([][]string, len(menu))
	for i := 1; i <= len(menu); i++ {
		column := make([]string, 4)
		column[0] = pterm.NewStyle(pterm.Bold, pterm.FgMagenta).Sprint(strconv.Itoa(i))
		column[1] = isVeg(menu[i].Type)
		column[2] = pterm.Italic.Sprint(menu[i].Name)
		column[3] = strconv.Itoa(menu[i].Cost)
		allRows = append(allRows, column)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"", "", "DISH", "(₹)"})
	table.SetHeaderColor(tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold})
	table.SetCenterSeparator("║")
	table.SetColumnSeparator("║")
	table.SetRowSeparator("═")
	table.AppendBulk(allRows)
	table.Render()
	fmt.Print("\n")
}
func wannaBuy() bool {
	var stopper bool
	pterm.DefaultCenter.Print(pterm.DefaultBox.Sprint(pterm.NewStyle(pterm.Bold, pterm.FgCyan).Sprint("Do you want to order anything from the menu? Answer with 'yes' or 'no', Enter 'exit' to quit")))
	for !stopper {
		choice := userInputHandler01()
		if choice == 0 {
			continue
		}
		switch choice {
		case 1:
			return true
		case 2:
			quit()
		}
	}
	return false
}

func userInputHandler02(lim int) int {
	pterm.FgMagenta.Print("~> ")
	choice := strings.ToLower(fetchVar01())
	switch choice {
	case "done":
		return 1
	case "exit":
		quit()
	case "":
		fmt.Print(pterm.Error.Sprint("Empty input!"), "\n\n")
		return 0
	default:
		intChoice, err := strconv.Atoi(choice)
		if err != nil {
			fmt.Print(pterm.Error.Sprint("Invalid input!"), "\n\n")
			return 0
		}
		if intChoice <= 0 {
			fmt.Print(pterm.Error.Sprint("Invalid input!"), "\n\n")
			return 0
		} else if intChoice > lim {
			fmt.Print(pterm.Error.Sprint("Input out of range!"), "\n\n")
			return 0
		} else {
			return intChoice + 1234
		}
	}
	return 0
}

func userInputHandler03() int {
	choice := fetchVar01()
	switch choice {
	case "":
		pterm.Error.Println("Empty input!")
		return 0
	default:
		intChoice, err := strconv.Atoi(choice)
		if err != nil {
			pterm.Error.Println("Invalid input!")
			return 0
		} else {
			if intChoice < 0 {
				pterm.Error.Println("Invalid input!")
				return 0
			} else {
				return intChoice + 1234
			}
		}
	}
}
func takeOrder(menu map[int]Menu) map[int]Bill {
	bill := make(map[int]Bill)
	var stopper bool
	fmt.Println(pterm.DefaultBox.WithTitle("Instructions: ").Sprintfln("* Enter %s to add the dish to the cart\n* Enter %s to checkout\n* Enter %s to quit", pterm.Bold.Sprint(pterm.FgMagenta.Sprint("DISH NUMBER")), pterm.Bold.Sprint(pterm.FgLightGreen.Sprint("DONE")), pterm.Bold.Sprint(pterm.FgRed.Sprint("EXIT"))))
	for !stopper {
		choice := userInputHandler02(len(menu))
		if choice == 0 {
			continue
		}
		switch choice {
		case 1:
			if len(bill) != 0 {
				stopper = true
			} else {
				fmt.Print(pterm.Error.Sprint("No Dishes ordered, instead enter 'exit' to quit"), "\n\n")
			}
		default:
			var howMany int
			choice -= 1234
			for {
				fmt.Printf("How many %s do you want ~> ", menu[choice].Name)
				howMany = userInputHandler03()
				if howMany != 0 {
					howMany -= 1234
					break
				}
			}
			var skipFlag, skipAddNewEntry bool
			for i := range bill {
				if bill[i].Name == menu[choice].Name {
					if howMany == 0 {
						pterm.Info.Printfln("%s was not added to cart", menu[choice].Name)
						pterm.Info.Printfln("%d %s already in cart\n", bill[choice].Qty, menu[choice].Name)
						skipFlag = true
						break
					}
					if entry, ok := bill[i]; ok {
						entry.Qty += howMany
						entry.Cost = entry.Qty * menu[choice].Cost
						bill[i] = entry
						pterm.FgLightGreen.Printfln("Added %d %s to cart, now %d in cart\n", howMany, menu[choice].Name, bill[i].Qty)
						skipAddNewEntry = true
					}
				}
			}
			if !skipFlag {
				if howMany == 0 {
					pterm.Info.Printfln("%s was not added to cart\n", menu[choice].Name)
					break
				}
				if !skipAddNewEntry {
					bill[len(bill)+1] = Bill{menu[choice].Type, menu[choice].Name, menu[choice].Cost * howMany, howMany}
					pterm.FgLightGreen.Printfln("Added %d %s to cart\n", howMany, menu[choice].Name)
				}
			}
		}
	}
	return bill
}
func printBill(bill map[int]Bill) int {
	allRows, total := make([][]string, len(bill)), 0
	for i := 1; i <= len(bill); i++ {
		column := make([]string, 5)
		column[0] = pterm.NewStyle(pterm.Bold, pterm.FgMagenta).Sprint(strconv.Itoa(i))
		column[1] = isVeg(bill[i].Type)
		column[2] = pterm.Italic.Sprint(bill[i].Name)
		column[3] = strconv.Itoa(bill[i].Qty)
		column[4] = strconv.Itoa(bill[i].Cost)
		allRows = append(allRows, column)
		total += bill[i].Cost
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"", "", "DISH", "Qty", "(₹)"})
	table.SetHeaderColor(tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold})
	table.SetCenterSeparator("║")
	table.SetColumnSeparator("║")
	table.SetRowSeparator("═")
	table.AppendBulk(allRows)
	fmt.Print("Your Bill:\n\n")
	table.Render()
	fmt.Print("\n")
	fmt.Print(pterm.DefaultBox.Sprintln("Total Cost:", pterm.NewStyle(pterm.Bold, pterm.Underscore, pterm.FgGreen).Sprint("₹ ", strconv.Itoa(total))))
	return total
}

func userInputHandler04(lim int) int {
	choice := fetchVar01()
	switch choice {
	case "":
		pterm.Error.Println("Empty input!")
		return 0
	default:
		intChoice, err := strconv.Atoi(choice)
		if err != nil {
			pterm.Error.Println("Invalid input!")
			return 0
		}
		if intChoice < 0 {
			pterm.Error.Println("Ivalid input!")
			return 0
		} else if intChoice > lim {
			pterm.Error.Println("Input out of range!")
			return 0
		} else {
			return intChoice + 1234
		}
	}
}
func wannaDel(bill map[int]Bill) (bool, map[int]Bill) {
	var stopper bool
	pterm.DefaultCenter.Print(pterm.DefaultBox.Sprint(pterm.NewStyle(pterm.Bold, pterm.FgCyan).Sprint("Do you want to remove any items? Answer with 'yes' or 'no', Enter 'exit' to quit")))
	for !stopper {
		choice := userInputHandler01()
		if choice == 0 {
			continue
		}
		switch choice {
		case 1:
			clearScr()
			printBill(bill)
			return false, letsDel(bill)
		case 2:
			return true, bill
		}
	}
	return false, bill
}
func letsDel(bill map[int]Bill) map[int]Bill {
	var stopper, wasRemoved bool
	var history []int
	billLen := len(bill)
	fmt.Print("\n")
	fmt.Println(pterm.DefaultBox.WithTitle("Instructions: ").Sprintfln("* Enter %s to remove the dish from the cart\n* Enter %s to save changes\n* Enter %s to quit", pterm.Bold.Sprint(pterm.FgMagenta.Sprint("DISH NUMBER")), pterm.Bold.Sprint(pterm.FgLightGreen.Sprint("DONE")), pterm.Bold.Sprint(pterm.FgRed.Sprint("EXIT"))))
	for !stopper {
		choice := userInputHandler02(billLen)
		if choice == 0 {
			continue
		}
		switch choice {
		case 1:
			if len(bill) != 0 {
				stopper = true
			} else {
				clearScr()
				pterm.DefaultCenter.Print(pterm.Italic.Sprint("😭️ All Items Removed from cart 😭️"))
				pterm.DefaultCenter.Print(pterm.Italic.Sprint("Thanks for visiting us, do come back later 😚️"))
				showCursor()
				os.Exit(0)
			}
		default:
			choice -= 1234
			var skipFlag bool
			for _, removedItem := range history {
				if choice == removedItem {
					fmt.Print(pterm.Error.Sprintfln("dish %d was already removed from bill, Enter 'done' to print the new bill", choice), "\n")
					skipFlag = true
				}
			}
			if !skipFlag {
				if bill[choice].Qty == 1 {
					pterm.FgLightGreen.Printfln("removed %s from cart\n", bill[choice].Name)
					history = append(history, choice)
					delete(bill, choice)
					wasRemoved = true
				} else {
					var stopper02 bool
					for !stopper02 {
						fmt.Printf("There are %d %s in cart, how many to remove ~> ", bill[choice].Qty, bill[choice].Name)
						howMany := userInputHandler04(bill[choice].Qty)
						if howMany == 0 {
							continue
						}
						stopper02 = true
						howMany -= 1234
						if howMany == 0 {
							pterm.Info.Printfln("Quantity of %s was not modified\n", bill[choice].Name)
						} else if howMany == bill[choice].Qty {
							pterm.FgLightGreen.Printfln("removed %s from cart\n", bill[choice].Name)
							history = append(history, choice)
							delete(bill, choice)
							wasRemoved = true
						} else {
							if entry, ok := bill[choice]; ok {
								entry.Cost -= (entry.Cost / entry.Qty) * howMany
								entry.Qty -= howMany
								bill[choice] = entry
								pterm.FgLightGreen.Printfln("%d %s now left in cart\n", bill[choice].Qty, bill[choice].Name)
							}
						}
					}
				}
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
	pterm.DefaultCenter.Print(pterm.DefaultBox.Sprint(pterm.NewStyle(pterm.Bold, pterm.FgCyan).Sprintf("Proceed to pay: %s, Confirm with 'yes' or 'no', Enter 'exit' to quit", pterm.NewStyle(pterm.FgGreen, pterm.Bold, pterm.Underscore).Sprint("₹ ", strconv.Itoa(price)))))
	for !stopper {
		choice := userInputHandler01()
		if choice == 0 {
			continue
		}
		switch choice {
		case 1:
			clearScr()
			pterm.Warning.Prefix = pterm.Prefix{
				Text:  "STATUS",
				Style: pterm.NewStyle(pterm.BgWhite, pterm.FgBlack),
			}
			hideCursor()
			pterm.Warning.Print("Processing Payment")
			time.Sleep(5 * time.Second)
			clearScr()
			hideCursor()
			pterm.Warning.Print("Almost Done")
			time.Sleep(3 * time.Second)
			clearScr()
			pterm.Success.Println("Payment Successful !")
			return true
		case 2:
			quit()
		}
	}
	return false
}
func showETA(bill map[int]Bill) {
	fmt.Printf("\nEnter your name (Optional) ~> ")
	username := fetchVar01()
	clearScr()
	if len(username) != 0 {
		pterm.Info.Printfln("Currently preparing for: %s", username)
	}
	var eta int
	for i := range bill {
		eta = eta + 4*bill[i].Qty // 4 seconds to make each dish
	}
	if eta >= 60 {
		pterm.Info.Printfln("ETA is %d min and %d sec, Order is:\n", eta/60, eta%60)
	} else {
		pterm.Info.Printfln("ETA is %d sec, Order is:\n", eta)
	}
}
func printOrder(bill map[int]Bill) {
	allRows := make([][]string, len(bill))
	for i := 1; i <= len(bill); i++ {
		column := make([]string, 3)
		column[0] = isVeg(bill[i].Type)
		column[1] = pterm.Italic.Sprint(bill[i].Name)
		column[2] = strconv.Itoa(bill[i].Qty)
		allRows = append(allRows, column)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"", "DISH", "Qty"})
	table.SetHeaderColor(tablewriter.Colors{}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold})
	table.SetCenterSeparator("║")
	table.SetColumnSeparator("║")
	table.SetRowSeparator("═")
	table.AppendBulk(allRows)
	table.Render()
	fmt.Print("\n")
}
func cookFood(bill map[int]Bill) {
	for i := range bill {
		fmt.Printf("Preparing %s now,\n", bill[i].Name)
		bar := progressbar.NewOptions(bill[i].Qty,
			progressbar.OptionSetRenderBlankState(true),
			progressbar.OptionFullWidth(),
			progressbar.OptionShowCount(),
			progressbar.OptionSetPredictTime(false),
			progressbar.OptionShowBytes(false),
		)
		for j := 1; j <= bill[i].Qty; j++ {
			time.Sleep(4 * time.Second)
			hideCursor()
			err := bar.Add(1)
			checkErr(err, "progress bar usage caused an error")
		}
		fmt.Print("\n")
		pterm.Info.Printfln("%s is ready !\n\n", bill[i].Name)
	}
}

func thankYou() {
	banner, err := (pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("THANK YOU !", pterm.NewStyle(pterm.FgMagenta))).Srender())
	checkErr(err, "pterm usage caused an error")
	fmt.Print("\n")
	pterm.DefaultCenter.Print(banner)
}
