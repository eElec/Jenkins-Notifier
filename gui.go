package main

import (
	"jenkins-notifier/utils"
	"jenkins-notifier/worker"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/getlantern/systray"
)

var (
	guiApp fyne.App
)

func runApp() {
	guiApp := app.New()
	w := guiApp.NewWindow("Hello")

	//* Menu
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Quit", func() { guiApp.Quit() }),
	)
	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", func() {
			dialog.ShowCustom("About", "Close", container.NewVBox(
				widget.NewLabel("Jenkins Notifier"),
				widget.NewLabel("Version: v0.1"),
				widget.NewLabel("Author: Adrish Aditya"),
			), w)
		}))

	mainMenu := fyne.NewMainMenu(
		fileMenu,
		helpMenu,
	)
	w.SetMainMenu(mainMenu)

	//* Layout

	w.Show()
	guiApp.Lifecycle().SetOnStopped(func() {
		log.Println("Quiting..")
		wg.Done()
		systray.Quit()
	})
	go systray.Run(onReady, onExit)
	guiApp.Run()
}

func onExit() {
	for k := range jobs {
		jobs[k].Stop()
		<-jobs[k].Event
	}
	if guiApp != nil {
		guiApp.Quit()
	}
}

func onReady() {
	systray.SetIcon(utils.GetIcon("icons/jenkins.ico"))
	systray.SetTitle("Jenkins Notifier")
	systray.SetTooltip("Jenkins Notifier")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuit.ClickedCh
		log.Println("Quiting...")
		systray.Quit()
	}()

	// job menu
	jobsMenu := systray.AddMenuItem("Jobs", "")

	for k := range jobs {
		subMenu := jobsMenu.AddSubMenuItemCheckbox(jobs[k].Name, "", jobs[k].IsRunning())

		//
		go func(subMenu *systray.MenuItem, job *worker.Job) {
			for {
				select {
				case <-subMenu.ClickedCh:
					job.TogglePause()
					if subMenu.Checked() {
						subMenu.Uncheck()
					} else {
						subMenu.Check()
					}
				}
			}
		}(subMenu, jobs[k])
	}
}
