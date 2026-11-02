package main

import (
	"bufio"
	"fmt"
	"ledger"
	"os"
	"strings"
)

func tryAdd(tx ledger.Transaction) {
	err := ledger.AddTransaction(tx)
	if err != nil {
		fmt.Println("Ошибка! Транзакция \""+tx.Description+"\" не добавлена:", err)
		return
	}
	fmt.Printf("Транзакция \"%s\" добавлена: %s, %.2f руб.\n", tx.Description, tx.Category, tx.Amount)
}

func trySetBudget(b ledger.Budget) {
	err := ledger.SetBudget(b)
	if err != nil {
		fmt.Println("Ошибка! Бюджет \""+b.Category+"\" не установлен:", err)
		return
	}
	fmt.Printf("Бюджет \"%s\" установлен: %.2f руб.\n", b.Category, b.Limit)
}

func printBudgets() {
	fmt.Println("Бюджеты:")
	for _, b := range ledger.ListBudgets() {
		fmt.Printf("  %s: %.2f руб. на %s\n", b.Category, b.Limit, b.Period)
	}
}

func main() {
	fmt.Println("Ledger service started")

	fmt.Println("\nПроверка через CheckValid:")
	items := []ledger.Validatable{
		ledger.Transaction{Amount: 500, Category: "еда", Date: "2026-11-01"},
		ledger.Transaction{Amount: 0, Category: "еда", Date: "2026-11-01"},
		ledger.Budget{Category: "еда", Limit: -100, Period: "месяц"},
	}
	for i, item := range items {
		err := ledger.CheckValid(item)
		if err != nil {
			fmt.Println("Проверка", i+1, "- ошибка:", err)
		} else {
			fmt.Println("Проверка", i+1, "- все в порядке")
		}
	}

	fmt.Println("\nУстанавливаем бюджеты:")
	trySetBudget(ledger.Budget{Category: "еда", Limit: 5000, Period: "месяц"})
	trySetBudget(ledger.Budget{Category: "транспорт", Limit: 2000, Period: "месяц"})
	trySetBudget(ledger.Budget{Category: "одежда", Limit: 0, Period: "месяц"})
	printBudgets()

	fmt.Println("\nТраты на еду (лимит 5000):")
	tryAdd(ledger.Transaction{Amount: 3000, Category: "еда", Description: "продукты на неделю", Date: "2026-11-01"})
	tryAdd(ledger.Transaction{Amount: 1500, Category: "еда", Description: "кафе", Date: "2026-11-02"})
	tryAdd(ledger.Transaction{Amount: 1000, Category: "еда", Description: "доставка пиццы", Date: "2026-11-03"})

	fmt.Println("\nКатегория без бюджета:")
	tryAdd(ledger.Transaction{Amount: 10000, Category: "подарки", Description: "подарок маме", Date: "2026-11-03"})

	fmt.Println("\nНеправильные транзакции:")
	tryAdd(ledger.Transaction{Amount: -200, Category: "еда", Description: "что-то странное", Date: "2026-11-03"})
	tryAdd(ledger.Transaction{Amount: 700, Category: "", Description: "непонятно что", Date: "2026-11-03"})

	fmt.Println("\nПоднимаем лимит на еду до 6000 и пробуем пиццу еще раз:")
	trySetBudget(ledger.Budget{Category: "еда", Limit: 6000, Period: "месяц"})
	tryAdd(ledger.Transaction{Amount: 1000, Category: "еда", Description: "доставка пиццы", Date: "2026-11-04"})

	fmt.Println("\nЗагружаем бюджеты из budgets.json")
	file, err := os.Open("budgets.json")
	if err != nil {
		fmt.Println("Ошибка! Не получилось открыть файл:", err)
	} else {
		defer file.Close()
		reader := bufio.NewReader(file)
		err = ledger.LoadBudgets(reader)
		if err != nil {
			fmt.Println("Ошибка!", err)
		} else {
			fmt.Println("Бюджеты загружены")
		}
	}
	printBudgets()

	fmt.Println("\nТраты после загрузки из файла:")
	tryAdd(ledger.Transaction{Amount: 2200, Category: "транспорт", Description: "проездной", Date: "2026-11-04"})
	tryAdd(ledger.Transaction{Amount: 2000, Category: "развлечения", Description: "концерт", Date: "2026-11-05"})
	tryAdd(ledger.Transaction{Amount: 1500, Category: "развлечения", Description: "кино и боулинг", Date: "2026-11-05"})

	fmt.Println("\nПробуем загрузить битый JSON (лимит строкой):")
	err = ledger.LoadBudgets(strings.NewReader(`[{"category": "еда", "limit": "много"}]`))
	if err != nil {
		fmt.Println("Ошибка!", err)
	}

	fmt.Println("\nСписок транзакций:")
	list := ledger.ListTransactions()
	for _, t := range list {
		fmt.Printf("  %d. %s, %.2f руб., %s (%s)\n", t.ID, t.Category, t.Amount, t.Description, t.Date)
	}
	fmt.Println("Всего транзакций:", len(list))
}
