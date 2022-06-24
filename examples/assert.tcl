//
// This demonstrates the use of the two words:
//
//  assert => Ensure that expr-results work as expected
//
//  assert_equal => Ensure that two things are identical.
//
// These are both implemented in our "standard library".


//
// Simple operations
//
assert 1        == 1
assert [+ 1 2]  == 3
assert [* 3 5]  == 15
assert 1        >= 1
assert 3        >= 1
assert 32       <= 32
assert 33       <= 1000
assert 31       ne "steve"
assert 32       ne 121
assert [/ 21 3] == 7


//
// String equality
//
set name "Steve"
assert_equal "Steve" "${name}"

//
// Numerical equality
//
assert_equal 343 [+ 340 3]
