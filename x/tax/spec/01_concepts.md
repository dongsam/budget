<!-- order: 1 -->

 # Concepts
## Tax Module

`x/tax` is a Cosmos SDK module that implements tax functionality that keeps track of the staking and provides tax rewards to farmers. One use case is to use this module to provide incentives for liquidity pool investors for their pool participation. 

## Taxes

There are two types of tax taxes in the `tax` module as below.

### 1. Public Tax Plan

A public tax plan can only be created through governance proposal meaning that the proporsal must be first agreed and passed in order to create a public plan.
### 2. Private Tax Plan

A private tax plan can be created with any account. The plan creator's account is used as distributing account `TaxPoolAddress` that will be distributed to farmers automatically. There is a fee `PlanCreationFee` paid upon plan creation to prevent from spamming attack. 

## Distribution Methods

There are two types of distribution methods  in the `tax` module as below.
### 1. Fixed Amount Plan

A `FixedAmountPlan` distributes fixed amount of coins to farmers for every epoch day. If the plan creators `TaxPoolAddress` is depleted with distributing coins, then there is no more coins to distribute unless it is filled up again.

### 2. Ratio Plan

A `RatioPlan` distributes to farmers by ratio distribution for every epoch day. If the plan creators `TaxPoolAddress` is depleted with distributing coins, then there is no more coins to distribute unless it is filled up with more coins.

