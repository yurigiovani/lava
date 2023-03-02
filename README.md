## Lava Network Test

Disclaimer
--
I've read that i should have use ignite scaffold, but i already have started 
using cosmos-sdk cloned from GitHub.

So, i've used only 2 modules beyond the module that i created, Bank and Auth.
* Bank module,  i've used to get balances from Module Account (Lottery Account) and
some AccAddress (clients).
* Account module, i've used to execute some automated tests for keeper that was created.

Steps
--
1. Execute the script
>./scripts/init.sh

The node will be started

2. After that execution, we can execute
> ./scripts/running-lottery.sh

It will be entered automatically on lottery for all clients that was created on step 1.

3. Once the network running, we can see the lottery commands:
> simd tx

> simd tx lottery


You'll see the lotteries command, on this case just ** lottery enter** command.

