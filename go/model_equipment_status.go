/*
 * 5G-EIR Equipment Identity Check
 *
 * 5G-EIR Equipment Identity Check Service
 *
 * API version: 1.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type EquipmentStatus string

// List of EquipmentStatus
const (
	WHITELISTED EquipmentStatus = "WHITELISTED"
	BLACKLISTED EquipmentStatus = "BLACKLISTED"
	GREYLISTED  EquipmentStatus = "GREYLISTED"
)
