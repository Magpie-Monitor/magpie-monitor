import { UrgencyLevel } from '@api/managment-service';
import PieChart from 'components/PieChart/PieChart';
import { PolarChartEntry } from 'components/PolarChart/PolarChart';
import colors from 'global/colors';

interface UrgancyPolarChartProps {
  urgencyIncidentCount: Record<UrgencyLevel, number>;
}

const urgencyToChartColor: Record<UrgencyLevel, string> = {
  LOW: colors.urgency.low,
  MEDIUM: colors.urgency.medium,
  HIGH: colors.urgency.high,
};

const CHART_COLOR_TRANSPARENCY = 'd0';

const UrgencyPolarChart = ({
  urgencyIncidentCount,
}: UrgancyPolarChartProps) => {
  const chartData: PolarChartEntry[] = Object.entries(urgencyIncidentCount).map(
    ([urgency, count]) => ({
      value: count,
      label: urgency,
      color:
        urgencyToChartColor[urgency as UrgencyLevel] + CHART_COLOR_TRANSPARENCY,
    }),
  );

  return <PieChart label="Incidents by Urgency" data={chartData} />;
};

export default UrgencyPolarChart;
