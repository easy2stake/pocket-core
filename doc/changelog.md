## RC-0.2.0
- Renamed RelayProof to Proof (in JSON)
- Renamed Invoice (memory) to Evidence
- Renamed StoredInvoice (blockchain persisted) to Receipt
- Renamed ProofWaitingPeriod to ClaimSubmissionWindow
- Changed RPC and from `node-proof` to `node-receipt`
- Update posmint module to use sdk.Ctx interface
- Fix `pseudorandomGenerator` unexported properties would return empty json
- Evidence now holds proof interface to allow for challenge proofs 
- Added Relay Request Hash (Hash of payload + meta) to RelayProof object
- Added Block to Dispatch Request
- Added Relay Meta field to relay request