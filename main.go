package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Barang struct {
	Id_barang   int    `json:"id_barang"`
	Nama_barang string `json:"nama_barang"`
	Deskripsi   string `json:"deskripsi"`
	Stok        int    `json:"stok"`
	Harga       int    `json:"harga"`
	Tgl_masuk   string `json:"tgl_masuk"`
}
type Transaksi struct {
	Id_transaksi int    `json:"id_transaksi"`
	Id_karyawan  int    `json:"karyawan"`
	Id_pembeli   int    `json:"id_pembeli"`
	Tgl_beli     string `json:"tgl_beli"`
	Id_barang    int    `json:"id_barang"`
	Jumlah_beli  int    `json:"jumlah_beli"`
	Harga        int    `json:"harga"`
}
type T struct {
	Id_transaksi int    `json:"id_transaksi"`
	Nama_kasir   string `json:"nama_kasir"`
	Id_pembeli   int    `json:"id_pembeli"`
	Nama_barang  string `json:"nama_barang"`
	Jumlah_beli  int    `json:"jumlah_beli"`
	Harga        int    `json:"harga"`
	Total        int    `json:"total"`
	Tgl_beli     string `json:"tgl_beli"`
}

func dtb() (db *sql.DB) {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/db_kasir")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Database Connected")
	}
	return db
}
func main() {

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{ //dari ini sampai 2 dbawah Optional
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// Route => handler Barang
	e.GET("/", HandlerIndex)
	//e.GET("/kd/:id_barang", HandlerGetBarang)
	e.POST("/tbarang", HandlerTbBarang)
	e.PUT("/ebarang/:id_barang", HandlerEBarang)
	e.DELETE("/dbarang/:id", HandlerDelBarang)

	//Route => Handler Transaksi
	e.POST("/ttrade", TbTrade)
	e.GET("/atrade", Atrade)
	e.GET("/trade/:id", tradeID)
	e.PUT("/etrade/:id", Etrade)
	e.DELETE("/dtrade/:id", Dtrade)

	e.Logger.Fatal(e.Start(":7000"))
}

func HandlerTbBarang(c echo.Context) error {

	trade := new(Barang)

	err1 := c.Bind(trade)
	if err1 != nil {
		return err1
	}
	db := dtb()
	q := "INSERT INTO barang(id_barang, nama_barang, stok, harga, tgl_masuk)VALUES(?, ?, ?, ?, ?)"
	ex, err2 := db.Prepare(q)

	if err2 != nil {
		fmt.Println(err2.Error())
	}
	defer ex.Close()
	hasil, err3 := ex.Exec(trade.Id_barang, trade.Nama_barang, trade.Stok, trade.Harga, trade.Tgl_masuk)

	if err3 != nil {
		panic(err3)
	}
	fmt.Println(hasil.LastInsertId())
	return c.JSON(http.StatusCreated, "Data "+trade.Nama_barang+" Berhasil Ditambahkan")
}

func HandlerEBarang(c echo.Context) error {
	idx := c.Param("id_barang")
	db := dtb()
	edb := new(Barang)
	errr := c.Bind(edb)
	if errr != nil {
		fmt.Println(errr)
	}

	q := "UPDATE barang SET nama_barang=? , deskripsi=? , stok=? , harga=? , tgl_masuk=? where id_barang=?"
	hasil, err1 := db.Query(q, edb.Nama_barang, edb.Deskripsi, edb.Stok, edb.Harga, edb.Tgl_masuk, idx)

	if err1 != nil {
		fmt.Println(err1)
	} else {
		fmt.Println(hasil)
		return c.JSON(http.StatusCreated, edb)
	}
	return c.JSON(http.StatusOK, edb.Id_barang)

}

func HandlerDelBarang(c echo.Context) error {
	idx := c.Param("id")
	db := dtb()

	q := "DELETE FROM barang where id_barang=?"
	hasil, err1 := db.Query(q, idx)

	if err1 != nil {
		fmt.Println(err1)
	} else {
		fmt.Println(hasil)
		return c.JSON(http.StatusOK, "Data Kode : "+idx+" Berhasil Dihapus")
	}
	return c.JSON(http.StatusOK, "Data kode: "+idx+" Telah diHapus")
}

func HandlerIndex(c echo.Context) error {
	db := dtb()
	dbs, errh := db.Query("select * from barang")
	if errh != nil {
		fmt.Println(errh)
	}

	brg := Barang{}
	allb := []Barang{}

	for dbs.Next() {
		var id_barang, stok, harga int
		var nama_barang, tgl_masuk, deskripsi string

		err1 := dbs.Scan(&id_barang, &nama_barang, &deskripsi, &stok, &harga, &tgl_masuk)
		if err1 != nil {
			fmt.Println(err1)
		}

		brg.Id_barang = id_barang
		brg.Nama_barang = nama_barang
		brg.Deskripsi = deskripsi
		brg.Stok = stok
		brg.Harga = harga
		brg.Tgl_masuk = tgl_masuk

		allb = append(allb, brg)
	}
	defer db.Close()
	return c.JSON(http.StatusOK, allb)

}

func TbTrade(c echo.Context) error {
	//idb := c.Param("id")
	db := dtb()

	tran := new(Transaksi)

	err1 := c.Bind(tran)
	if err1 != nil {
		return err1
	}

	q := "INSERT INTO transaksi(id_transaksi,id_karyawan, id_pembeli, tgl_beli)VALUES(?, ?, ?, ?)"
	q1 := "INSERT INTO detail_transaksi(id_barang,id_transaksi,jumlah_beli)VALUES(?,?,?)"
	exe, err2 := db.Prepare(q)
	exez, errj := db.Prepare(q1)

	if err2 != nil {
		fmt.Println(err2.Error())
	}
	if errj != nil {
		fmt.Println(errj.Error())
	}
	defer exe.Close()

	hasil, err3 := exe.Exec(tran.Id_transaksi, tran.Id_karyawan, tran.Id_pembeli, tran.Tgl_beli)
	hasil2, err4 := exez.Exec(tran.Id_barang, tran.Id_transaksi, tran.Jumlah_beli)

	if err3 != nil {
		panic(err3)
	}
	if err4 != nil {
		panic(err4)
	}

	fmt.Println(hasil.LastInsertId())
	fmt.Println(hasil2.LastInsertId())
	return c.JSON(http.StatusCreated, "Data Transaksi Berhasil Ditambahkan")

}

func Atrade(c echo.Context) error {
	db := dtb()
	dbb, errg := db.Query("select  transaksi.id_transaksi,karyawan.nama,transaksi.id_pembeli,barang.nama_barang,detail_transaksi.jumlah_beli,barang.harga,transaksi.tgl_beli,detail_transaksi.id_barang,barang.id_barang,transaksi.id_karyawan,karyawan.id_karyawan from detail_transaksi,karyawan,barang,transaksi where detail_transaksi.id_transaksi=transaksi.id_transaksi AND detail_transaksi.id_barang=barang.id_barang AND transaksi.id_karyawan=karyawan.id_karyawan")
	if errg != nil {
		fmt.Println(errg)
	}

	tr := T{}
	traa := []T{}

	for dbb.Next() {
		var id_transaksi, jumlah_beli, harga, id_karyawan, id_barang, id_pembeli int
		var nama_barang, tgl_beli, nama string

		err2 := dbb.Scan(&id_transaksi, &nama, &id_pembeli, &nama_barang, &jumlah_beli, &harga, &tgl_beli, &id_karyawan, &id_barang, &id_barang, &id_karyawan)
		if err2 != nil {
			fmt.Println(err2)
		}

		tr.Id_transaksi = id_transaksi
		tr.Nama_kasir = nama
		tr.Id_pembeli = id_pembeli
		tr.Nama_barang = nama_barang
		tr.Jumlah_beli = jumlah_beli
		tr.Harga = harga
		tr.Total = harga * jumlah_beli
		tr.Tgl_beli = tgl_beli

		traa = append(traa, tr)
	}
	defer db.Close()
	return c.JSON(http.StatusOK, traa)

}

func tradeID(c echo.Context) error {
	idb := c.Param("id")
	db := dtb()
	dbb, errg := db.Query("select  transaksi.id_transaksi,karyawan.nama,transaksi.id_pembeli,barang.nama_barang,detail_transaksi.jumlah_beli,barang.harga,transaksi.tgl_beli,detail_transaksi.id_barang,barang.id_barang,transaksi.id_karyawan,karyawan.id_karyawan from detail_transaksi,karyawan,barang,transaksi where detail_transaksi.id_transaksi=transaksi.id_transaksi AND detail_transaksi.id_barang=barang.id_barang AND transaksi.id_karyawan=karyawan.id_karyawan AND detail_transaksi.id_transaksi=?", idb)
	if errg != nil {
		fmt.Println(errg)
	}

	tr := T{}
	traa := []T{}

	for dbb.Next() {
		var id_transaksi, jumlah_beli, harga, id_karyawan, id_barang, id_pembeli int
		var nama_barang, tgl_beli, nama string

		err2 := dbb.Scan(&id_transaksi, &nama, &id_pembeli, &nama_barang, &jumlah_beli, &harga, &tgl_beli, &id_karyawan, &id_barang, &id_barang, &id_karyawan)
		if err2 != nil {
			fmt.Println(err2)
		}

		tr.Id_transaksi = id_transaksi
		tr.Nama_kasir = nama
		tr.Id_pembeli = id_pembeli
		tr.Nama_barang = nama_barang
		tr.Jumlah_beli = jumlah_beli
		tr.Harga = harga
		tr.Total = harga * jumlah_beli
		tr.Tgl_beli = tgl_beli

		traa = append(traa, tr)
	}
	defer db.Close()
	return c.JSON(http.StatusOK, traa)

}

func Etrade(c echo.Context) error {
	idb := c.Param("id")
	db := dtb()

	tr := new(Transaksi)
	err1 := c.Bind(tr)
	if err1 != nil {
		fmt.Println(err1)
	}
	qq := "UPDATE detail_transaksi SET id_barang=?, jumlah_beli=? where id_detail=?"
	result, err2 := db.Query(qq, tr.Id_barang, tr.Jumlah_beli, idb)
	if err2 != nil {
		fmt.Println(err2)
	} else {
		fmt.Println(result)
		return c.JSON(http.StatusCreated, tr.Id_barang)
	}
	return c.JSON(http.StatusOK, "data Berhasil dirubah kode detail"+idb)

}

func Dtrade(c echo.Context) error {
	idx := c.Param("id")
	db := dtb()

	q := "DELETE FROM transaksi where id_transaksi=?"
	hasil, err1 := db.Query(q, idx)

	if err1 != nil {
		fmt.Println(err1)
	} else {
		fmt.Println(hasil)
		return c.JSON(http.StatusOK, "Data Kode : "+idx+" Berhasil Dihapus")
	}
	return c.JSON(http.StatusOK, "Data kode: "+idx+" Telah diHapus")
}
