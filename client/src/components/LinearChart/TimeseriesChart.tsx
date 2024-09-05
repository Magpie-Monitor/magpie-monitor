import { Line } from "react-chartjs-2";
import "chartjs-adapter-moment";
import { useEffect, useState } from "react";
import { ChartData, ChartOptions } from "chart.js";
import "./TimeseriesChart.scss";

interface TimeSeriesChartProps {
  data: TimeSeriesData[];
  legend: string;
  title: string;
  from?: number; // Unix timestamp
  to?: number; // Unix timestamp
  color: string;
  backgroundColor: string;
}

export interface TimeSeriesData {
  timestamp: number; // Unix timestamp
  value: number;
}

export const getChartOptions = (
  title: string,
  // from?: number,
  // to?: number,
): ChartOptions<"line"> => ({
  responsive: true,
  maintainAspectRatio: true,
  resizeDelay: 200,
  color: "white",
  backgroundColor: "none",
  scales: {
    y: {
      beginAtZero: true,
      ticks: {
        color: "white",
      },
    },
    x: {
      type: "time",
      time: {
        unit: "hour",
        displayFormats: {
          hour: "MMM D hA",
        },
      },
      ticks: {
        color: "white",
        maxTicksLimit: 5,
      },
      // min: from,
      // max: to,
    },
  },
  plugins: {
    legend: {
      display: true,
      position: "bottom",
    },
    title: {
      display: true,
      text: title,
      position: "top",
      color: "white",
      font: {
        size: 16,
      },
    },
  },
});

export const TimeSeriesChart = ({
  data,
  from,
  to,
  legend,
  title,
  color,
  backgroundColor,
}: TimeSeriesChartProps) => {
  const [chartData, setChartData] = useState<ChartData<"line"> | null>(null);

  useEffect(() => {
    setChartData({
      labels: [],
      datasets: [
        {
          label: legend,
          data: data
            .sort((a, b) => a.timestamp - b.timestamp)
            .map((d) => ({
              x: d.timestamp,
              y: d.value,
            })),
          fill: true,
          borderColor: color,
          backgroundColor: backgroundColor,
          tension: 0.1,
        },
      ],
    });
  }, [data, from, to, legend, title, color, backgroundColor]);

  const timeSeriesChartOptions = getChartOptions(title);

  return (
    <div className="timeseries-chart">
      {chartData && <Line data={chartData} options={timeSeriesChartOptions} />}
    </div>
  );
};
