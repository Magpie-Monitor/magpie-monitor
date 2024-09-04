import { Chart as ChartJS, ChartOptions, RadialLinearScale } from "chart.js";
import { Radar } from "react-chartjs-2";
import "./RadarChart.scss";

ChartJS.register(RadialLinearScale);

export interface RadarChartEntry {
  value: number;
  label: string;
  color: string;
}

export interface RadarChartProps {
  data: RadarChartEntry[];
  label: string;
}

const RadarChart = ({ data, label }: RadarChartProps) => {
  const config = {
    labels: data.map((entry) => entry.label),
    datasets: [
      {
        label: label,
        data: data.map((entry) => entry.value),
        backgroundColor: data.map((entry) => entry.color),
        borderColor: "transparent",
        borderWidth: 1,
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: true,
    resizeDelay: 200,
    color: "white",
    backgroundColor: "none",
    elements: {
      line: {
        borderWidth: 3,
        backgroundColor: "rgba(75, 192, 192, 0.6)",
        borderColor: "rgba(75, 192, 192, 1)",
      },
    },
    scales: {
      r: {
        grid: {
          color: "white",
        },
        angleLines: {
          color: "white",
        },
        pointLabels: {
          color: "white",
          font: {
            size: 14,
            color: "white",
          },
        },
        ticks: {
          color: "white",
          backdropColor: "black",
          font: {
            size: 18,
            weight: "bold",
          },
        },
      },
    },

    plugins: {
      legend: {
        display: true,
        position: "bottom",
      },
      title: {
        display: true,
        text: label,
        position: "top",
        color: "white",
        font: {
          size: 16,
        },
      },
    },
  } as ChartOptions<"radar">;

  return (
    <div className="polar-chart">
      <Radar data={config} options={options} />
    </div>
  );
};

export default RadarChart;
