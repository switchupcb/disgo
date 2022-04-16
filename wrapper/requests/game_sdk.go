package requests

// Get Achievements
// GET /applications/{application.id}/get-set-up}/achievements
// TODO
type GetAchievements struct {
	// TODO
}

// Get Achievement
// GET /applications/{application.id}/get-set-up}/achievements/{achievement.id}/data-models-achievement-struct}
// TODO
type GetAchievement struct {
	// TODO
}

// Create Achievement
// POST /applications/{application.id}/get-set-up}/achievements
// TODO
type CreateAchievement struct {
	// TODO
}

// Update Achievement
// PATCH /applications/{application.id}/get-set-up}/achievements/{achievement.id}/data-models-achievement-struct}
// TODO
type UpdateAchievement struct {
	// TODO
}

// Delete Achievement
// DELETE /applications/{application.id}/get-set-up}/achievements/{achievement.id}/data-models-achievement-struct}
// TODO
type DeleteAchievement struct {
	// TODO
}

// Update User Achievement
// PUT /applications/{application.id}/get-set-up}/achievements/{achievement.id}/data-models-achievement-struct}
// TODO
type UpdateUserAchievement struct {
	// TODO
}

// Get User Achievements
// GET /users/@me/applications/{application.id}/get-set-up}/achievements
// TODO
type GetUserAchievements struct {
	// TODO
}

/// .go filegame_sdk\Lobbies.md
// Create Lobby
// POST /lobbies
// TODO
type CreateLobby struct {
	// TODO
}

// Update Lobby
// PATCH /lobbies/{lobby.id}/data-models-lobby-struct}
// TODO
type UpdateLobby struct {
	// TODO
}

// Delete Lobby
// DELETE /lobbies/{lobby.id}/data-models-lobby-struct}
// TODO
type DeleteLobby struct {
	// TODO
}

// Update Lobby Member
// PATCH /lobbies/{lobby.id}/data-models-lobby-struct}/members/{user.id}/user-object}
// TODO
type UpdateLobbyMember struct {
	// TODO
}

// Create Lobby Search
// POST /lobbies/search
// TODO
type CreateLobbySearch struct {
	// TODO
}

// Send Lobby Data
// POST /lobbies/{lobby.id}/data-models-lobby-struct}/send
// TODO
type SendLobbyData struct {
	// TODO
}

/// .go filegame_sdk\Store.md
// Get Entitlements
// GET /applications/{application.id}/get-set-up}/entitlements
// TODO
type GetEntitlements struct {
	// TODO
}

// Get Entitlement
// GET /applications/{application.id}/get-set-up}/entitlements/{entitlement.id}/data-models-entitlement-struct}
// TODO
type GetEntitlement struct {
	// TODO
}

// Get SKUs
// GET /applications/{application.id}/get-set-up}/skus
// TODO
type GetSKUs struct {
	// TODO
}

// Consume SKU
// POST /applications/{application.id}/get-set-up}/entitlements/{entitlement.id}/data-models-entitlement-struct}/consume
// TODO
type ConsumeSKU struct {
	// TODO
}

// Delete Test Entitlement
// DELETE /applications/{application.id}/get-set-up}/entitlements/{entitlement.id}/data-models-entitlement-struct}
// TODO
type DeleteTestEntitlement struct {
	// TODO
}

// Create Purchase Discount
// PUT /store/skus/{sku.id}/data-models-sku-struct}/discounts/{user.id}/user-object}
// TODO
type CreatePurchaseDiscount struct {
	// TODO
}

// Delete Purchase Discount
// DELETE /store/skus/{sku.id}/data-models-sku-struct}/discounts/{user.id}/user-object}
// TODO
type DeletePurchaseDiscount struct {
	// TODO
}
