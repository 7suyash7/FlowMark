transaction(publicKey: String, numOfKeysToAdd: Int) {
  prepare(signer: AuthAccount) {
    let bytes = publicKey.decodeHex()
    let key = PublicKey(
      publicKey: bytes,
      signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
    )

    var counter = 0
    while counter < numOfKeysToAdd {
      counter = counter + 1
      signer.keys.add(
        publicKey: key,
        hashAlgorithm: HashAlgorithm.SHA3_256,
        weight: 0.0
      )
    }
  }

  execute {

  }
}