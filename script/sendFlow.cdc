import FungibleToken from 0xee82856bf20e2aa6
import FlowToken from 0x0ae53cb6e3f42a79

transaction(amount: UFix64, recipient: Address) {
    let sentVault: @FungibleToken.Vault
    prepare(signer: AuthAccount) {
    let vaultRef = signer.borrow<&FlowToken.Vault>(from: /storage/flowTokenVault)
        ?? panic("failed to borrow reference to sender vault")

    self.sentVault <- vaultRef.withdraw(amount: amount)
    }

    execute {
    let receiverRef =  getAccount(recipient)
        .getCapability(/public/flowTokenReceiver)
        .borrow<&{FungibleToken.Receiver}>()
        ?? panic("failed to borrow reference to recipient vault")

    receiverRef.deposit(from: <-self.sentVault)
    }
}