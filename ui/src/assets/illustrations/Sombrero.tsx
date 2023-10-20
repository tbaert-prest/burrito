import { SVGProps } from "react";

const Sombrero = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={424}
    height={424}
    viewBox="0 0 424 424"
    fill="none"
    {...props}
  >
    <path
      fill="#F463A0"
      stroke="#000"
      strokeWidth={5}
      d="M261.5 244h3.045l-.593-2.986-19.039-95.988C241.798 129.317 228.015 118 212 118c-16.015 0-29.798 11.317-32.913 27.026l-19.039 95.988-.593 2.986H261.5ZM410 246.5V244H14v2.5c0 8.696 5.863 16.762 15.943 23.857 10.14 7.138 24.833 13.517 43.118 18.848 18.313 5.34 40.004 9.561 63.825 12.438C160.711 304.52 186.234 306 212 306c25.766 0 51.289-1.48 75.114-4.357 23.821-2.877 45.512-7.098 63.825-12.438 18.285-5.331 32.978-11.71 43.118-18.848C404.137 263.262 410 255.196 410 246.5Z"
    />
    <mask
      id="a-sombrero"
      width={392}
      height={58}
      x={16}
      y={246}
      maskUnits="userSpaceOnUse"
      style={{
        maskType: "alpha",
      }}
    >
      <path
        fill="#F463A0"
        d="M407.5 246.5c0 7.485-5.057 14.897-14.882 21.813-9.824 6.916-24.225 13.199-42.379 18.492-18.153 5.293-39.705 9.492-63.424 12.356-23.72 2.865-49.142 4.339-74.815 4.339-25.673 0-51.095-1.474-74.815-4.339-23.719-2.864-45.27-7.063-63.424-12.356-18.154-5.293-32.555-11.576-42.38-18.492C21.558 261.397 16.5 253.985 16.5 246.5h391Z"
      />
    </mask>
    <g fill="#6BEF70" stroke="#000" strokeWidth={5} mask="url(#a-sombrero)">
      <path d="m17.33 268.74 2.17 3.799 2.17-3.799 12-21 2.138-3.74H3.192l2.137 3.74 12 21ZM52.33 268.74l2.17 3.799 2.17-3.799 12-21 2.138-3.74H38.192l2.137 3.74 12 21ZM87.33 268.74l2.17 3.799 2.17-3.799 12.001-21 2.137-3.74H73.192l2.137 3.74 12 21ZM122.329 268.74l2.171 3.799 2.171-3.799 12-21 2.137-3.74h-32.616l2.137 3.74 12 21ZM157.329 268.74l2.171 3.799 2.171-3.799 12-21 2.137-3.74h-32.616l2.137 3.74 12 21ZM192.329 268.74l2.171 3.799 2.171-3.799 12-21 2.137-3.74h-32.616l2.137 3.74 12 21ZM227.329 268.74l2.171 3.799 2.171-3.799 12-21 2.137-3.74h-32.616l2.137 3.74 12 21ZM262.329 268.74l2.171 3.799 2.171-3.799 12-21 2.137-3.74h-32.616l2.137 3.74 12 21ZM297.329 268.74l2.171 3.799 2.171-3.799 12-21 2.137-3.74h-32.616l2.137 3.74 12 21ZM332.329 268.74l2.171 3.799 2.171-3.799 12-21 2.137-3.74h-32.616l2.137 3.74 12 21ZM367.329 268.74l2.171 3.799 2.171-3.799 12-21 2.137-3.74h-32.616l2.137 3.74 12 21ZM402.329 268.74l2.171 3.799 2.171-3.799 12-21 2.137-3.74h-32.616l2.137 3.74 12 21Z" />
    </g>
  </svg>
);

export default Sombrero;
