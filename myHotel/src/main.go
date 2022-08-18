package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	mergeMaps "github.com/imdario/mergo"
	"github.com/olekukonko/tablewriter"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/schollz/progressbar/v3"
)

func main() {
	preCheck()
	clearScr()
	greet()
	menu := decideMenu()
	clearScr()
	if printMenu(menu); wannaBuy() {
		clearScr()
		printMenu(menu)
		bill := takeOrder(menu)
		clearScr()
		printBill(bill)
		var isItSame bool
		for {
			isItSame, bill = wannaDel(bill)
			clearScr()
			if isItSame {
				if printBill(bill); askToPay() {
					username := showETA(findETA(bill))
					logOrder(bill, username)
					printOrder(bill)
					cookFood(bill)
					showCursor()
					break
				}
			} else {
				printBill(bill)
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

func preCheck() {
	if runtime.GOOS != "linux" {
		clearScr()
		fmt.Println("Unsupported Platform:", runtime.GOOS)
		fmt.Println("Supported Platform: linux only with kitty terminal")
		showCursor()
		os.Exit(0)
	}
}
func hideCursor() {
	fmt.Print("\033[?25l")
}
func showCursor() {
	fmt.Print("\033[?25h")
}
func clearScr() {
	fmt.Print("\x1bc")
	hideCursor()
}
func quit() {
	clearScr()
	pterm.DefaultCenter.Print(pterm.Italic.Sprint("Thanks for visiting us, do come back later 😚️"))
	showCursor()
	os.Exit(0)
}
func checkErr(err error, hint string) {
	if err != nil {
		clearScr()
		fmt.Println("Oops, " + hint + " usage caused an error :(\n")
		showCursor()
		panic(err)
	}
}
func fetchVar() string {
	reader := bufio.NewReader(os.Stdin)
	hideCursor()
	input, err := reader.ReadString('\n')
	checkErr(err, "fetchVar()-->bufio")
	return strings.TrimSpace(input)
}
func notifyErr(message string) {
	cmd := exec.Command("notify-send", message)
	checkErr(cmd.Run(), "notifyErr()-->cmd.Run()")
}
func isItVeg(checkDish string) string {
	switch checkDish {
	case "veg":
		return "🟢️"
	case "nonveg":
		return "🔴️"
	default:
		return "  "
	}
}

func greet() {
	banner, err := pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("W", pterm.NewStyle(pterm.FgLightRed)),
		putils.LettersFromStringWithStyle("E", pterm.NewStyle(pterm.FgLightYellow)),
		putils.LettersFromStringWithStyle("L", pterm.NewStyle(pterm.FgLightRed)),
		putils.LettersFromStringWithStyle("C", pterm.NewStyle(pterm.FgLightYellow)),
		putils.LettersFromStringWithStyle("O", pterm.NewStyle(pterm.FgLightRed)),
		putils.LettersFromStringWithStyle("M", pterm.NewStyle(pterm.FgLightYellow)),
		putils.LettersFromStringWithStyle("E", pterm.NewStyle(pterm.FgLightRed)),
		putils.LettersFromStringWithStyle(" !", pterm.NewStyle(pterm.FgLightMagenta))).Srender()
	checkErr(err, "greet()-->pterm+putils")
	pterm.DefaultCenter.Println(banner)
	pterm.DefaultCenter.Printfln("🎊️ %s Exclusive Offers 🎊️", time.Now().Format("Monday"))
	pterm.DefaultCenter.Printfln("Enjoy discounts on all dishes every %s !", time.Now().Format("Monday"))
}

func userInputHandler01() int {
	pterm.NewStyle(pterm.Bold, pterm.FgLightMagenta).Print("~~> ")
	choice := strings.ToLower(fetchVar())
	switch choice {
	case "yes":
		return 1
	case "no":
		return 2
	case "exit":
		quit()
	case "":
		fmt.Print(pterm.Error.Sprint("Empty input!"), "\n\n")
		notifyErr("Empty input!")
		return 0
	default:
		fmt.Print(pterm.Error.Sprint("Invalid input!"), "\n\n")
		notifyErr("Invalid input!")
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
	vegMenu := map[int]Menu{
		1: {"veg", "Burger", 250},
		2: {"veg", "Fries", 100},
		3: {"veg", "Pizza", 350},
		4: {"veg", "Sandwich", 150},
		5: {"veg", "Soup", 200},
		6: {"veg", "Noodles", 350},
		7: {"veg", "Curry", 300},
		8: {"veg", "Rice", 250},
		// 9: {"veg", "Dish", Cost},
	}
	return vegMenu
}
func setMenuNonVeg() map[int]Menu {
	vegMenu := setMenuVeg()
	vegLen := len(vegMenu)
	nonVegMenu := map[int]Menu{
		vegLen + 1: {"nonveg", "Chicken Burger", 300},
		vegLen + 2: {"nonveg", "Chicken Wings", 200},
		vegLen + 3: {"nonveg", "Chicken Pizza", 400},
		vegLen + 4: {"nonveg", "Chicken Sandwich", 200},
		vegLen + 5: {"nonveg", "Chicken Curry", 350},
		// vegLen + 6: {"nonveg", "Dish", Cost},
	}
	err := mergeMaps.Merge(&nonVegMenu, vegMenu)
	checkErr(err, "setMenuNonVeg()-->mergeMaps")
	return nonVegMenu
}

func printMenu(menu map[int]Menu) {
	allRows := make([][]string, 0, len(menu))
	for i := 1; i <= len(menu); i++ {
		column := make([]string, 4)
		column[0] = pterm.NewStyle(pterm.Bold, pterm.FgLightMagenta).Sprint(strconv.Itoa(i))
		column[1] = isItVeg(menu[i].Type)
		column[2] = pterm.Italic.Sprint(menu[i].Name)
		column[3] = strconv.Itoa(menu[i].Cost)
		allRows = append(allRows, column)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"", "", "DISH", "₹"})
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
	pterm.NewStyle(pterm.Bold, pterm.FgLightMagenta).Print("~~> ")
	switch choice := strings.ToLower(fetchVar()); choice {
	case "done":
		return 1
	case "exit":
		quit()
	case "":
		fmt.Print(pterm.Error.Sprint("Empty input!"), "\n\n")
		notifyErr("Empty input!")
		return 0
	default:
		intChoice, err := strconv.Atoi(choice)
		if err != nil || intChoice <= 0 {
			fmt.Print(pterm.Error.Sprint("Invalid input!"), "\n\n")
			notifyErr("Invalid input!")
			return 0
		} else if intChoice > lim {
			fmt.Print(pterm.Error.Sprint("Input out of range!"), "\n\n")
			notifyErr("Input out of range!")
			return 0
		}
		return intChoice + 1234
	}
	return 0
}

func userInputHandler03() int {
	switch choice := fetchVar(); choice {
	case "":
		pterm.Error.Println("Empty input!")
		notifyErr("Empty input!")
		return 0
	default:
		intChoice, err := strconv.Atoi(choice)
		if err != nil || intChoice < 0 {
			pterm.Error.Println("Invalid input!")
			notifyErr("Invalid input!")
			return 0
		}
		return intChoice + 1234
	}
}
func takeOrder(menu map[int]Menu) map[int]Bill {
	bill := make(map[int]Bill)
	fmt.Println(pterm.DefaultBox.WithTitle("Instructions: ").Sprintfln("* Enter %s to add the dish to the cart\n* Enter %s to checkout\n* Enter %s to quit", pterm.Bold.Sprint(pterm.FgLightMagenta.Sprint("DISH NUMBER")), pterm.Bold.Sprint(pterm.FgGreen.Sprint("DONE")), pterm.Bold.Sprint(pterm.FgRed.Sprint("EXIT"))))
	var stopper bool
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
				notifyErr("No Dishes ordered, instead enter 'exit' to quit")
			}
		default:
			choice -= 1234
			var howMany int
			for {
				fmt.Printf("How many %s do you want ~~> ", menu[choice].Name)
				howMany = userInputHandler03()
				if howMany != 0 {
					howMany -= 1234
					break
				}
			}
			var skipFlag bool
			for i := range bill {
				if bill[i].Name == menu[choice].Name {
					if howMany == 0 {
						pterm.Info.Printfln("%s was not added to cart, already %d in cart\n", menu[choice].Name, bill[choice].Qty)
						skipFlag = true
						break
					}
					if entry, ok := bill[i]; ok {
						entry.Qty += howMany
						entry.Cost = entry.Qty * menu[choice].Cost
						bill[i] = entry
						pterm.FgGreen.Printfln("Added %d %s to cart, now %d in cart\n", howMany, menu[choice].Name, bill[i].Qty)
						skipFlag = true
						break
					}
				}
			}
			if !skipFlag && howMany != 0 {
				bill[len(bill)+1] = Bill{menu[choice].Type, menu[choice].Name, menu[choice].Cost * howMany, howMany}
				pterm.FgGreen.Printfln("Added %d %s to cart\n", howMany, menu[choice].Name)
			} else if !skipFlag && howMany == 0 {
				pterm.Info.Printfln("%s was not added to cart\n", menu[choice].Name)
			}
		}
	}
	return bill
}

func printBill(bill map[int]Bill) {
	allRows, total := make([][]string, 0, len(bill)), 0
	for i := 1; i <= len(bill); i++ {
		column := make([]string, 5)
		column[0] = pterm.NewStyle(pterm.Bold, pterm.FgLightMagenta).Sprint(strconv.Itoa(i))
		column[1] = isItVeg(bill[i].Type)
		column[2] = pterm.Italic.Sprint(bill[i].Name)
		column[3] = "x " + strconv.Itoa(bill[i].Qty)
		column[4] = strconv.Itoa(bill[i].Cost)
		allRows = append(allRows, column)
		total += bill[i].Cost
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"", "", "DISH", "Qty", "₹"})
	table.SetHeaderColor(tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold})
	table.SetCenterSeparator("║")
	table.SetColumnSeparator("║")
	table.SetRowSeparator("═")
	table.AppendBulk(allRows)
	table.Render()
	pterm.DefaultCenter.Print(pterm.DefaultBox.Sprint("    Total Cost: ", pterm.NewStyle(pterm.Bold, pterm.Underscore, pterm.FgLightGreen).Sprint("₹ ", strconv.Itoa(total)), "    "))
}

func userInputHandler04(lim int) int {
	switch choice := fetchVar(); choice {
	case "":
		pterm.Error.Println("Empty input!")
		notifyErr("Empty input!")
		return 0
	default:
		intChoice, err := strconv.Atoi(choice)
		if err != nil || intChoice < 0 {
			pterm.Error.Println("Invalid input!")
			notifyErr("Invalid input!")
			return 0
		} else if intChoice > lim {
			pterm.Error.Println("Input out of range!")
			notifyErr("Input out of range!")
			return 0
		}
		return intChoice + 1234
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
	return false, nil
}
func letsDel(bill map[int]Bill) map[int]Bill {
	var stopper, wasRemoved bool
	var history []int
	billLen := len(bill)
	fmt.Print("\n")
	fmt.Println(pterm.DefaultBox.WithTitle("Instructions: ").Sprintfln("* Enter %s to remove the dish from the cart\n* Enter %s to save changes\n* Enter %s to quit", pterm.Bold.Sprint(pterm.FgLightMagenta.Sprint("DISH NUMBER")), pterm.Bold.Sprint(pterm.FgGreen.Sprint("DONE")), pterm.Bold.Sprint(pterm.FgRed.Sprint("EXIT"))))
	for !stopper {
		choice := userInputHandler02(billLen)
		if choice == 0 {
			continue
		}
	switchBlock:
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
			for _, removedItem := range history {
				if choice == removedItem {
					fmt.Print(pterm.Error.Sprintfln("dish %d was already removed from bill, Enter 'done' to print the new bill", choice), "\n")
					notifyErr(fmt.Sprintf("dish %d was already removed from bill, Enter 'done' to print the new bill", choice))
					break switchBlock
				}
			}
			if bill[choice].Qty == 1 {
				pterm.FgGreen.Printfln("removed %s from cart\n", bill[choice].Name)
				history = append(history, choice)
				delete(bill, choice)
				wasRemoved = true
			} else {
				var stopper02 bool
				for !stopper02 {
					fmt.Printf("There are %d %s in cart, how many to remove ~~> ", bill[choice].Qty, bill[choice].Name)
					howMany := userInputHandler04(bill[choice].Qty)
					if howMany == 0 {
						continue
					}
					stopper02 = true
					howMany -= 1234
					switch howMany {
					case 0:
						pterm.Info.Printfln("Quantity of %s was not modified\n", bill[choice].Name)
					case bill[choice].Qty:
						pterm.FgGreen.Printfln("removed %s from cart\n", bill[choice].Name)
						history = append(history, choice)
						delete(bill, choice)
						wasRemoved = true
					default:
						if entry, ok := bill[choice]; ok {
							entry.Cost -= (entry.Cost / entry.Qty) * howMany
							entry.Qty -= howMany
							bill[choice] = entry
							pterm.FgGreen.Printfln("%d %s now left in cart\n", bill[choice].Qty, bill[choice].Name)
						}
					}
				}
			}
		}
	}
	if wasRemoved {
		return rearrangeMapKeys(bill)
	}
	return bill
}
func rearrangeMapKeys(bill map[int]Bill) map[int]Bill {
	newBill, ctr := make(map[int]Bill), 1
	for key, val := range bill {
		newBill[ctr] = bill[key]
		newBill[ctr] = val
		ctr++
	}
	return newBill
}

func askToPay() bool {
	var stopper bool
	pterm.DefaultCenter.Print(pterm.DefaultBox.Sprint(pterm.NewStyle(pterm.Bold, pterm.FgCyan).Sprintf("Proceed to pay? Confirm with 'yes' or 'no', Enter 'exit' to quit")))
	for !stopper {
		choice := userInputHandler01()
		if choice == 0 {
			continue
		}
		switch choice {
		case 1:
			transact()
			return true
		case 2:
			quit()
		}
	}
	return false
}
func transact() {
	pterm.Warning.Prefix = pterm.Prefix{
		Text:  "STATUS",
		Style: pterm.NewStyle(pterm.BgWhite, pterm.FgBlack),
	}
	clearScr()
	pterm.Warning.Print("Processing Payment...")
	time.Sleep(4 * time.Second)
	clearScr()
	pterm.Warning.Print("Almost Done")
	time.Sleep(2 * time.Second)
	clearScr()
	pterm.Success.Println("Payment Successful !")
}

func findETA(bill map[int]Bill) (int, string) {
	fmt.Printf("\nEnter your name (Optional) ~~> ")
	username := fetchVar()
	clearScr()
	var eta int
	for i := range bill {
		eta += 4 * bill[i].Qty // 4 seconds to make each dish
	}
	return eta, username
}
func showETA(eta int, username string) string {
	if len(username) != 0 {
		pterm.Info.Printfln("👨‍🍳️is currently preparing for: %s", username)
	}
	if eta < 60 {
		pterm.Info.Printfln("ETA is %d sec, Order is:\n", eta)
	} else {
		pterm.Info.Printfln("ETA is %d min and %d sec, Order is:\n", eta/60, eta%60)
	}
	return username
}

func printOrder(bill map[int]Bill) {
	allRows := make([][]string, 0, len(bill))
	for i := 1; i <= len(bill); i++ {
		column := make([]string, 3)
		column[0] = isItVeg(bill[i].Type)
		column[1] = pterm.Italic.Sprint(bill[i].Name)
		column[2] = "x " + strconv.Itoa(bill[i].Qty)
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

func logOrder(bill map[int]Bill, username string) {
	var fileName string
	timeFormat := time.Now().Format("January-2__Monday__3:4:5_PM__Username[")
	if len(username) != 0 {
		fileName = timeFormat + username + "].txt"
	} else {
		fileName = timeFormat + "Unknown" + "].txt"
	}
	var err error
	if err = os.Chdir("Order_History"); os.IsNotExist(err) {
		checkErr(os.Mkdir("Order_History", os.ModePerm), "logOrder()-->os.Mkdir()")
		err = os.Chdir("Order_History")
	}
	checkErr(err, "logOrder()-->os.Chdir()")
	file, err := os.Create(fileName)
	checkErr(err, "logOrder()-->os.Create()")
	_, err = fmt.Fprint(file, bill)
	checkErr(err, "logOrder()-->fmt.Fprint()")
	checkErr(file.Close(), "logOrder()-->file.Close()")
	checkErr(os.Chdir(".."), "logOrder()-->os.Chdir()")
}

func cookFood(bill map[int]Bill) {
	var totalQty int
	for i := range bill {
		totalQty += bill[i].Qty
	}
	bar := progressbar.NewOptions(totalQty,
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionShowBytes(false),
		progressbar.OptionOnCompletion(thankYou),
	)
	for i := range bill {
		for j := 1; j <= bill[i].Qty; j++ {
			hideCursor()
			time.Sleep(4 * time.Second)
			err := bar.Add(1)
			checkErr(err, "cookFood()-->progressbar")
		}
	}
}

func thankYou() {
	banner, err := pterm.DefaultBigText.WithLetters(putils.LettersFromStringWithStyle("THANK YOU !", pterm.NewStyle(pterm.FgLightMagenta))).Srender()
	checkErr(err, "thankYou()-->pterm+putils")
	fmt.Print("\n\n\n\n")
	pterm.DefaultCenter.Print(banner)
}
