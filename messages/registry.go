// Code generated by go GenerateMessage; DO NOT EDIT.
// This file was generated by robots at
// 2021-08-30 22:10:37.7887453 +0200 CEST m=+57.965421001

package messages

var Types_ = map[int]Message{
	CharacterSelectionID:                          &CharacterSelection{},
	HelloConnectID:                                &HelloConnect{},
	IdentificationID:                              &Identification{},
	SelectedServerDataExtendedID:                  &SelectedServerDataExtended{},
	AbstractCharacterInformationID:                &abstractCharacterInformation{},
	AccountCapabilitiesID:                         &accountCapabilities{},
	AccountHouseID:                                &accountHouse{},
	AccountHouseInformationsID:                    &accountHouseInformations{},
	AccountTagInformationID:                       &accountTagInformation{},
	AchievementAchievedID:                         &achievementAchieved{},
	AchievementAchievedRewardableID:               &achievementAchievedRewardable{},
	AchievementListID:                             &achievementList{},
	AcquaintancesGetListID:                        &acquaintancesGetList{},
	AlignmentRankUpdateID:                         &alignmentRankUpdate{},
	AllianceInformationID:                         &allianceInformation{},
	AllianceInsiderPrismInformationID:             &allianceInsiderPrismInformation{},
	AlliancePrismInformationID:                    &alliancePrismInformation{},
	AlmanachCalendarDateID:                        &almanachCalendarDate{},
	ArenaLeagueRankingID:                          &arenaLeagueRanking{},
	ArenaRankInfosID:                              &arenaRankInfos{},
	ArenaRankingID:                                &arenaRanking{},
	AuthenticationTicketID:                        &authenticationTicket{},
	AuthenticationTicketAcceptedID:                &authenticationTicketAccepted{},
	BasicAllianceInformationsID:                   &basicAllianceInformations{},
	BasicCharactersListID:                         &basicCharactersList{},
	BasicGuildInformationsID:                      &basicGuildInformations{},
	BasicNamedAllianceInformationID:               &basicNamedAllianceInformation{},
	BasicNoOperationID:                            &basicNoOperation{},
	BasicPongID:                                   &basicPong{},
	BasicTimeID:                                   &basicTime{},
	CharacterBaseInformationsID:                   &characterBaseInformations{},
	CharacterBasicMinimalInformationsID:           &characterBasicMinimalInformations{},
	CharacterCapabilitiesID:                       &characterCapabilities{},
	CharacterCharacteristicForPresetID:            &characterCharacteristicForPreset{},
	CharacterExperienceGainID:                     &characterExperienceGain{},
	CharacterLoadingCompleteID:                    &characterLoadingComplete{},
	CharacterMinimalInformationsID:                &characterMinimalInformations{},
	CharacterMinimalPlusLookInformationsID:        &characterMinimalPlusLookInformations{},
	CharacterSelectedSuccessID:                    &characterSelectedSuccess{},
	CharactersListID:                              &charactersList{},
	CharactersListRequestID:                       &charactersListRequest{},
	ChatAbstractServerID:                          &chatAbstractServer{},
	ChatCommunityChannelCommunityID:               &chatCommunityChannelCommunity{},
	ChatServerID:                                  &chatServer{},
	CheckIntegrityID:                              &checkIntegrity{},
	ClientKeyID:                                   &clientKey{},
	CredentialsAcknowledgementID:                  &credentialsAcknowledgement{},
	EmoteAddID:                                    &emoteAdd{},
	EmoteListID:                                   &emoteList{},
	EnabledChannelsID:                             &enabledChannels{},
	EntitiesPresetID:                              &entitiesPreset{},
	EntityLookID:                                  &entityLook{},
	ForgettableSpellsPresetID:                     &forgettableSpellsPreset{},
	FriendGuildWarnOnAchievementCompleteStateID:   &friendGuildWarnOnAchievementCompleteState{},
	FriendStatusShareStateID:                      &friendStatusShareState{},
	FriendWarnOnConnectionStateID:                 &friendWarnOnConnectionState{},
	FriendWarnOnLevelGainStateID:                  &friendWarnOnLevelGainState{},
	FriendsGetListID:                              &friendsGetList{},
	FullStatsPresetID:                             &fullStatsPreset{},
	GameContextCreateRequestID:                    &gameContextCreateRequest{},
	GameRolePlayArenaUpdatePlayerInfosID:          &gameRolePlayArenaUpdatePlayerInfos{},
	GameRolePlayArenaUpdatePlayerInfosAllQueuesID: &gameRolePlayArenaUpdatePlayerInfosAllQueues{},
	GameServerInformationsID:                      &gameServerInformations{},
	GuildEmblemID:                                 &guildEmblem{},
	GuildInformationsID:                           &guildInformations{},
	GuildMemberWarnOnConnectionStateID:            &guildMemberWarnOnConnectionState{},
	HaapiApiKeyID:                                 &haapiApiKey{},
	HaapiApiKeyRequestID:                          &haapiApiKeyRequest{},
	HaapiSessionID:                                &haapiSession{},
	HavenBagPackListMessageID:                     &havenBagPackListMessage{},
	HavenBagRoomPreviewInformationID:              &havenBagRoomPreviewInformation{},
	HavenBagRoomUpdateID:                          &havenBagRoomUpdate{},
	HelloGameID:                                   &helloGame{},
	HouseGuildedInformationsID:                    &houseGuildedInformations{},
	HouseInformationsID:                           &houseInformations{},
	HouseInformationsForGuildID:                   &houseInformationsForGuild{},
	HouseInformationsInsideID:                     &houseInformationsInside{},
	HouseInstanceInformationsID:                   &houseInstanceInformations{},
	HouseOnMapInformationsID:                      &houseOnMapInformations{},
	IconNamedPresetID:                             &iconNamedPreset{},
	IdentificationFailedID:                        &identificationFailed{},
	IdentificationFailedForBadVersionID:           &identificationFailedForBadVersion{},
	IdentificationSuccessID:                       &identificationSuccess{},
	IdolID:                                        &idol{},
	IdolListID:                                    &idolList{},
	IdolsPresetID:                                 &idolsPreset{},
	IgnoredGetListID:                              &ignoredGetList{},
	InteractiveElementNamedSkillID:                &interactiveElementNamedSkill{},
	InteractiveElementSkillID:                     &interactiveElementSkill{},
	InventoryContentID:                            &inventoryContent{},
	InventoryWeightID:                             &inventoryWeight{},
	ItemID:                                        &item{},
	ItemForPresetID:                               &itemForPreset{},
	ItemWrapperID:                                 &itemWrapper{},
	ItemsPresetID:                                 &itemsPreset{},
	JobCrafterDirectorySettingsID:                 &jobCrafterDirectorySettings{},
	JobCrafterDirectorySettingsMessageID:          &jobCrafterDirectorySettingsMessage{},
	JobDescriptionID:                              &jobDescription{},
	JobDescriptionMID:                             &jobDescriptionM{},
	JobExperienceID:                               &jobExperience{},
	JobExperienceMultiUpdateID:                    &jobExperienceMultiUpdate{},
	KnownZaapListID:                               &knownZaapList{},
	LoginQueueStatusID:                            &loginQueueStatus{},
	MountClientID:                                 &mountClient{},
	MountSetID:                                    &mountSet{},
	MountXpRatioID:                                &mountXpRatio{},
	NotificationListID:                            &notificationList{},
	ObjectEffectID:                                &objectEffect{},
	ObjectEffectCreatureID:                        &objectEffectCreature{},
	ObjectEffectDateID:                            &objectEffectDate{},
	ObjectEffectDiceID:                            &objectEffectDice{},
	ObjectEffectDurationID:                        &objectEffectDuration{},
	ObjectEffectIntegerID:                         &objectEffectInteger{},
	ObjectEffectLadderID:                          &objectEffectLadder{},
	ObjectEffectMinMaxID:                          &objectEffectMinMax{},
	ObjectEffectMountID:                           &objectEffectMount{},
	ObjectEffectStringID:                          &objectEffectString{},
	ObjectItemID:                                  &objectItem{},
	ObjectItemInformationWithQuantityID:           &objectItemInformationWithQuantity{},
	PartyIdolID:                                   &partyIdol{},
	PresetID:                                      &preset{},
	PresetsID:                                     &presets{},
	PresetsContainerPresetID:                      &presetsContainerPreset{},
	PrismGeolocalizedInformationID:                &prismGeolocalizedInformation{},
	PrismInformationID:                            &prismInformation{},
	PrismSubareaEmptyInfoID:                       &prismSubareaEmptyInfo{},
	PrismsListID:                                  &prismsList{},
	ProtocolID:                                    &protocol{},
	RawDataID:                                     &rawData{},
	SelectedServerDataID:                          &selectedServerData{},
	SequenceNumberRequestID:                       &sequenceNumberRequest{},
	ServerExperienceModificatorID:                 &serverExperienceModificator{},
	ServerOptionalFeaturesID:                      &serverOptionalFeatures{},
	ServerSessionConstantID:                       &serverSessionConstant{},
	ServerSessionConstantIntegerID:                &serverSessionConstantInteger{},
	ServerSessionConstantLongID:                   &serverSessionConstantLong{},
	ServerSessionConstantStringID:                 &serverSessionConstantString{},
	ServerSessionConstantsID:                      &serverSessionConstants{},
	ServerSettingsID:                              &serverSettings{},
	SetUpdateID:                                   &setUpdate{},
	ShortcutID:                                    &shortcut{},
	ShortcutBarContentID:                          &shortcutBarContent{},
	ShortcutEmoteID:                               &shortcutEmote{},
	ShortcutEntitiesPresetID:                      &shortcutEntitiesPreset{},
	ShortcutObjectID:                              &shortcutObject{},
	ShortcutObjectIdolsPresetID:                   &shortcutObjectIdolsPreset{},
	ShortcutObjectItemID:                          &shortcutObjectItem{},
	ShortcutObjectPresetID:                        &shortcutObjectPreset{},
	ShortcutSmileyID:                              &shortcutSmiley{},
	ShortcutSpellID:                               &shortcutSpell{},
	SimpleCharacterCharacteristicForPresetID:      &simpleCharacterCharacteristicForPreset{},
	SkillActionDescriptionID:                      &skillActionDescription{},
	SkillActionDescriptionCollectID:               &skillActionDescriptionCollect{},
	SkillActionDescriptionCraftID:                 &skillActionDescriptionCraft{},
	SkillActionDescriptionTimedID:                 &skillActionDescriptionTimed{},
	SpellForPresetID:                              &spellForPreset{},
	SpellItemID:                                   &spellItem{},
	SpellListID:                                   &spellList{},
	SpellsPresetID:                                &spellsPreset{},
	SpouseGetInformationsID:                       &spouseGetInformations{},
	SpouseStatusID:                                &spouseStatus{},
	StartupActionAddObjectID:                      &startupActionAddObject{},
	StatsPresetID:                                 &statsPreset{},
	StorageKamasUpdateID:                          &storageKamasUpdate{},
	SubEntityID:                                   &subEntity{},
	TextInformationID:                             &textInformation{},
	TrustCertificateID:                            &trustCertificate{},
	TrustStatusID:                                 &trustStatus{},
	VersionID:                                     &version{},
	WarnOnPermaDeathStateID:                       &warnOnPermaDeathState{},
}
