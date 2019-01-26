package word

import (
	"fmt"
	"math/rand"
	"time"
	"unicode"
	"unicode/utf8"
)

import "testing"

func TestIsPalindrome(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"", true},
		{"a", true},
		{"aa", true},
		{"ab", false},
		{"kayak", true},
		{"detartrated", true},
		{"A man, a plan, a canal: Panama", true},
		{"Evil I did dwell; lewd did I live.", true},
		{"Able was I ere I saw Elba", true},
		{"été", true},
		{"Et se resservir, ivresse reste.", true},
		{"palindrome", false}, 
		{"desserts", false},   
	}
	for _, test := range tests {
		if got := IsPalindrome(test.input); got != test.want {
			t.Errorf("IsPalindrome(%q) = %v", test.input, got)
		}
	}
}

func BenchmarkIsPalindrome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPalindrome("A man, a plan, a canal: Panama")
	}
}

func ExampleIsPalindrome() {
	fmt.Println(IsPalindrome("A man, a plan, a canal: Panama"))
	fmt.Println(IsPalindrome("palindrome"))
	

	
}


func randomPalindrome(rng *rand.Rand) string {
	n := rng.Intn(25) 
	runes := make([]rune, n)
	for i := 0; i < (n+1)/2; i++ {
		r := rune(rng.Intn(0x1000)) 
		runes[i] = r
		runes[n-1-i] = r
	}
	return string(runes)
}

func randomNotPalindrome(rng *rand.Rand) string {
	n := rng.Intn(25) 
	runes := make([]rune, n)
	for i := 0; i < (n+1)/2; i++ {
		for {
			c := rng.Intn(0x999)
			r := rune(c) 
			r2 := rune(c + 1)
			if unicode.IsLetter(r) == true && unicode.IsLetter(r2) == true && unicode.ToLower(r) != unicode.ToLower(r2) {
				runes[i] = r
				runes[n-1-i] = r2
				break
			}
		}

	}
	return string(runes)
}

func TestRandomPalindromes(t *testing.T) {
	
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 1000; i++ {
		p := randomPalindrome(rng)
		if !IsPalindrome(p) {
			t.Errorf("IsPalindrome(%q) = false", p)
		}
	}
}

func TestRandomNotPalindromes(t *testing.T) {
	
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 1000; i++ {
		p := randomNotPalindrome(rng)
		if utf8.RuneCountInString(p) > 1 && IsPalindrome(p) {
			t.Errorf("IsPalindrome(%q) = true", p)
		}
	}
}

func TestSingleNotPalidromes(t *testing.T) {
	var p string
	p = "Ĕĕ"
	if !IsPalindrome(p) {
		t.Errorf("IsPalindrome(%q) = false", p)
	}
}


