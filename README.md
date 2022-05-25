### Purpose
This repo has 2 main goals:
1) Demonstrate that logrus is quite a slow library in comparison to zap
   1) simply run `go test -bench . -run=XXX`
2) Demonstrate the impossibility (at the moment, 2022/05/25) for logrus to properly log the caller line if it's being wrapped
   1) check `Test_LogrusCallerWithWrapper` vs `Test_ZapCallerWithWrapper`

There is also an initial Logger interface implementation and a zap implementation.