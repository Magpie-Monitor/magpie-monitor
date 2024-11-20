package pl.pwr.zpi.reports.enums;

public enum ReportType {
    ON_DEMAND,
    SCHEDULED;
    public static ReportType fromString(String value) {
        if (value == null) {
            throw new IllegalArgumentException("Value cannot be null");
        }

        return switch (value.trim().toLowerCase()) {
            case "on-demand", "on_demand", "ondemand" -> ON_DEMAND;
            case "scheduled" -> SCHEDULED;
            default -> throw new IllegalArgumentException("Unknown ReportType: " + value);
        };
    }
}