package migrate

import "github.com/jinzhu/gorm"

//Migrate performs basic migration for test environment
func Migrate(db *gorm.DB) {
	db.DropTable(&CalcStepsLog{})
	db.DropTable(&CalcStepsLogDefinition{})
	db.DropTable(&RateDefinition{})
	db.DropTable(&RateFormula{})
	db.DropTable(&Rate{})
	db.DropTable(&Release{})
	db.DropTable(&WorksWriter{})
	db.DropTable(&Writer{})
	db.DropTable(&Resource{})
	db.DropTable(&UsageSummary{})
	db.DropTable(&Header{})
	db.DropTable(&Service{})
	db.DropTable(&Work{})

	db.CreateTable(&Header{})
	db.CreateTable(&Release{})
	db.CreateTable(&Resource{})
	db.CreateTable(&Work{})
	db.CreateTable(&Writer{})
	db.CreateTable(&WorksWriter{})

	db.Model(Resource{}).AddForeignKey("origin_id", "usage.headers(header_id)", "CASCADE", "CASCADE")
	db.Model(Resource{}).AddForeignKey("work_id", "usage.works(work_id)", "CASCADE", "CASCADE")
	db.Model(WorksWriter{}).AddForeignKey("work_id", "usage.works(work_id)", "CASCADE", "CASCADE")
	db.Model(WorksWriter{}).AddForeignKey("party_id", "usage.writers(party_id)", "CASCADE", "CASCADE")

	db.CreateTable(&RateFormula{})
	db.CreateTable(&Service{})
	db.CreateTable(&UsageSummary{})
	db.CreateTable(&RateDefinition{})
	db.CreateTable(&Rate{})
	db.CreateTable(&CalcStepsLogDefinition{})
	db.CreateTable(&CalcStepsLog{})

	db.Model(UsageSummary{}).AddForeignKey("service_id", "royalty.services(service_id)", "CASCADE", "CASCADE")
	db.Model(UsageSummary{}).AddForeignKey("header_id", "usage.headers(header_id)", "CASCADE", "CASCADE")
	db.Model(RateDefinition{}).AddForeignKey("service_id", "royalty.services(service_id)", "CASCADE", "CASCADE")
	db.Model(RateDefinition{}).AddForeignKey("rate_formula_id", "royalty.rate_formulas(rate_formula_id)", "CASCADE", "CASCADE")
	db.Model(Rate{}).AddForeignKey("usage_summary_id", "usage.usage_summaries(usage_summary_id)", "CASCADE", "CASCADE")
	db.Model(CalcStepsLogDefinition{}).AddForeignKey("service_id", "royalty.services(service_id)", "CASCADE", "CASCADE")
	db.Model(CalcStepsLog{}).AddForeignKey("usage_summary_id", "usage.usage_summaries(usage_summary_id)", "CASCADE", "CASCADE")
	db.Model(CalcStepsLog{}).AddForeignKey("log_definition_id", "royalty.calc_steps_log_definitions(log_definition_id)", "CASCADE",
		"CASCADE")
	orderedHexUUID(db)
	orderedBinUUID(db)
	unorderedUUID(db)

}

func orderedHexUUID(db *gorm.DB) {
	orderedHexUUID := `
		create function usage.ordered_hex_uuid(uuid bytea) returns text
        immutable
    	language plpgsql
		as
		$$
		BEGIN
    		RETURN pg_catalog.encode(uuid, 'hex');
		END ;
		$$;`

	db.Exec(
		orderedHexUUID)
}

func unorderedUUID(db *gorm.DB) {
	unorderedUUID := `
	create function usage.unordered_uuid(encoded_uuid bytea) returns uuid
	immutable
	language plpgsql
	as
	$$
	DECLARE
	uuid TEXT;
	BEGIN
	uuid := usage.ordered_hex_uuid(encoded_uuid);
	RETURN concat_ws('-', substring(uuid FROM 9 FOR 8), substring(uuid FROM 5 FOR 4), substring(uuid FROM 1 FOR 4),
	substring(uuid FROM 17 FOR 4), substring(uuid FROM 21))::UUID;
	END ;
	$$;`

	db.Exec(
		unorderedUUID)
}

func orderedBinUUID(db *gorm.DB) {
	orderedBinUUID := `
	create function usage.ordered_bin_uuid(text_uuid text) returns bytea
    immutable
    language plpgsql
	as
	$$
	BEGIN
    RETURN pg_catalog.decode(
                  substring(text_uuid FROM 15 FOR 4) || substring(text_uuid FROM 10 FOR 4) ||
                  substring(text_uuid FROM 1 FOR 8) ||
                  substring(text_uuid FROM 20 FOR 4) || substring(text_uuid FROM 25),
                  'hex');
	END ;
	$$;`

	db.Exec(
		orderedBinUUID)

}
