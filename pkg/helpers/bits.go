package helpers

// Установка бита
func Bit_set[T uint8](n T, pos uint8) T {
  return n | (1 << pos)
}

// Сброс бита
func Bit_clear[T uint8](n T, pos uint8) T {
  return n &^ (1 << pos)
}

// Проверка бита
func Bit_has[T uint8](n T, pos uint8) bool {
  return (n & (1 << pos)) > 0
}

// Установка бита
func BitMask_set[T uint8](n T, bm T) T {
  return n | bm
}

// Сброс бита
func BitMask_clear[T uint8](n T, bm T) T {
  return n &^ bm
}

// Проверка бита
func BitMask_has[T uint8](n T, bm T) bool {
  return (n & bm) == bm
}
