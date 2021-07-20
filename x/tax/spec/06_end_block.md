<!-- order: 6 -->

 # End-Block

- Termination of Tax Plan
    - Private Plan
        - distribution stops
        - remove plan states
        - keep stake, reward states for unstakable stakes and claimable rewards each farmers
        - rest of the fund in `taxPoolAddress` sent to `terminationAddress`, but in Private Plan case, `taxPoolAddress` == `terminationAddress`, so the fund is not moved
    - Public Plan
        - distribution stops
        - remove plan states
        - keep stake, reward states for unstakable stakes and claimable rewards each farmers
        - rest of the fund in `taxPoolAddress` sent to `terminationAddress`