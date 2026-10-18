package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Inisialisasi Aplikasi dan Jendela
	myApp := app.New()
	myWindow := myApp.NewWindow("Contoh Aplikasi Chat Sederhana")

	// Area untuk menampilkan riwayat chat
	chatArea := widget.NewMultiLineEntry()
	chatArea.Wrapping = fyne.TextWrapWord
	chatArea.SetText("Sistem: Selamat datang di aplikasi chat contoh.\n")
	
	// Membuat area chat tidak bisa diedit oleh user (hanya untuk tampilan)
	chatArea.Disable() 

	// Kolom input untuk mengetik pesan
	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("Ketik pesan Anda di sini...")

	// Tombol Kirim
	sendBtn := widget.NewButton("Kirim", func() {
		pesan := inputEntry.Text
		if pesan != "" {
			// Tambahkan pesan ke area chat
			chatArea.Enable() // Enable sementara untuk mengubah teks
			chatArea.SetText(chatArea.Text + "Anda: " + pesan + "\n")
			chatArea.Disable() // Disable kembali
			
			// Kosongkan kolom input
			inputEntry.SetText("")
		}
	})

	// Mengatur layout input dan tombol (Tombol di sebelah kanan input)
	inputBox := container.NewBorder(nil, nil, nil, sendBtn, inputEntry)

	// Menggabungkan area chat dan kotak input secara vertikal
	content := container.NewVBox(
		chatArea,
		inputBox,
	)

	// Set konten ke jendela dan jalankan
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 500))
	myWindow.ShowAndRun()
}
