import type { Config } from 'jest';

const config: Config = {
  preset: 'jest-preset-angular',
  setupFilesAfterEnv: ['<rootDir>/src/setup-jest.ts'],
  testEnvironment: 'jsdom',
  transform: {
    '^.+\\.(ts|mjs|js|html)$': [
      'jest-preset-angular',
      {
        tsconfig: '<rootDir>/tsconfig.jest.json',
        stringifyContentPathRegex: '\\.(html|svg)$',
      },
    ],
  },
  moduleFileExtensions: ['ts', 'html', 'js', 'json', 'mjs'],
  testMatch: ['**/src/**/*.spec.ts'],
  coverageDirectory: 'coverage/jest',
  collectCoverageFrom: [
    'src/app/**/*.ts',
    '!src/app/**/*.module.ts',
    '!src/main.ts',
  ],
};

export default config;