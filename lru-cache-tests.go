package main

import "fmt"

func TestScenario1_FirstAddition() {
	fmt.Println("\n### Test Scenario 1: First Addition (Empty Cache) ###")
	cache := NewLRUCache[string](2)

	fmt.Println("Initial state:")
	cache.Print()

	fmt.Println("\nAdding key 'A' with value 'ValueA'")
	cache.Add("A", "ValueA")

	cache.Print()
	fmt.Println("Expected: Head=A, Tail=A, CurrentSize=1")
}

func TestScenario2_SecondAddition() {
	fmt.Println("\n### Test Scenario 2: Second Addition ###")
	cache := NewLRUCache[string](2)
	cache.Add("A", "ValueA")

	fmt.Println("Initial state (after first add):")
	cache.Print()

	fmt.Println("\nAdding key 'B' with value 'ValueB'")
	cache.Add("B", "ValueB")

	cache.Print()
	fmt.Println("Expected: Head=B->A, Tail=A, CurrentSize=2")
}

func TestScenario3_EvictionWhenFull() {
	fmt.Println("\n### Test Scenario 3: Eviction When Full ###")
	cache := NewLRUCache[string](2)
	cache.Add("A", "ValueA")
	cache.Add("B", "ValueB")

	fmt.Println("Initial state (cache full):")
	cache.Print()

	fmt.Println("\nAdding key 'C' with value 'ValueC' (should evict A)")
	cache.Add("C", "ValueC")

	cache.Print()
	fmt.Println("Expected: Head=C->B, Tail=B, CurrentSize=2, A should be evicted")
}

func TestScenario4_UpdateExistingKey() {
	fmt.Println("\n### Test Scenario 4: Update Existing Key ###")
	cache := NewLRUCache[string](2)
	cache.Add("A", "ValueA")
	cache.Add("B", "ValueB")

	fmt.Println("Initial state:")
	cache.Print()

	fmt.Println("\nUpdating key 'A' with value 'NewValueA' (should move to head)")
	cache.Add("A", "NewValueA")

	cache.Print()
	fmt.Println("Expected: Head=A->B, Tail=B, CurrentSize=2 (no eviction)")
}

func TestScenario5_LookupMovesToHead() {
	fmt.Println("\n### Test Scenario 5: Lookup Moves to Head ###")
	cache := NewLRUCache[string](2)
	cache.Add("A", "ValueA")
	cache.Add("B", "ValueB")

	fmt.Println("Initial state:")
	cache.Print()

	fmt.Println("\nLooking up key 'A' (should move to head)")
	val, exists := cache.Lookup("A")
	fmt.Printf("Lookup result: value=%v, exists=%v\n", val, exists)

	cache.Print()
	fmt.Println("Expected: Head=A->B, Tail=B, CurrentSize=2")
}

func TestScenario6_ComplexSequence() {
	fmt.Println("\n### Test Scenario 6: Complex Sequence ###")
	cache := NewLRUCache[int](3)

	operations := []struct {
		op    string
		key   string
		value int
	}{
		{"Add", "A", 1},
		{"Add", "B", 2},
		{"Add", "C", 3},
		{"Lookup", "A", 0},
		{"Add", "D", 4},
		{"Update", "B", 20},
	}

	for i, op := range operations {
		fmt.Printf("\n--- Step %d: %s(%s) ---\n", i+1, op.op, op.key)

		switch op.op {
		case "Add", "Update":
			cache.Add(op.key, op.value)
		case "Lookup":
			val, exists := cache.Lookup(op.key)
			fmt.Printf("Lookup result: value=%v, exists=%v\n", val, exists)
		}

		cache.Print()
	}

	fmt.Println("\nFinal expected state: Head=B->D->A, Tail=A, CurrentSize=3")
	fmt.Println("(C was evicted when D was added, since B and A were more recently used)")
}
