package g_learning_connector

func PanicIfNeeded(err error) {
	if err != nil {
		panic(err)
	}
}
