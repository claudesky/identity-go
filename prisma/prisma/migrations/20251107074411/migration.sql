/*
  Warnings:

  - You are about to drop the column `client_id` on the `clients` table. All the data in the column will be lost.

*/
-- DropForeignKey
ALTER TABLE "public"."authorization_codes" DROP CONSTRAINT "authorization_codes_client_id_fkey";

-- DropIndex
DROP INDEX "public"."clients_client_id_key";

-- AlterTable
ALTER TABLE "clients" DROP COLUMN "client_id";

-- AddForeignKey
ALTER TABLE "authorization_codes" ADD CONSTRAINT "authorization_codes_client_id_fkey" FOREIGN KEY ("client_id") REFERENCES "clients"("id") ON DELETE CASCADE ON UPDATE NO ACTION;
